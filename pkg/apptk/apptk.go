package apptk

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
)

// AppToolkit 提供给Application的工具包
type AppToolkit struct {
}

// ReadIntParam 获取id并转换uint64
func (app *AppToolkit) ReadIntParam(r *http.Request, field string) (uint64, *gotk.ApiError) {
	val := chi.URLParam(r, field)
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil || id < 1 {
		errVal := fmt.Errorf("无效的 %s 参数: %s", field, val)
		return 0, errs.ErrBadRequest.With(errVal, errVal.Error())
	}
	return uint64(id), nil
}

// ReadPageQuery 从url获取pageNum、pageSize、sortFields三个字段
func (app *AppToolkit) ReadPageQuery(r *http.Request) dbrepo.Filters {
	var filter dbrepo.Filters
	pageNumText := r.URL.Query().Get("pageNum")
	pageSizeText := r.URL.Query().Get("pageSize")

	// 无需处理错误，0值使用时，会设置为默认值
	filter.PageNum, _ = strconv.Atoi(pageNumText)
	filter.PageSize, _ = strconv.Atoi(pageSizeText)

	// 需要支持两种传参数方式：仅仅适合少量数组元素传递，url长度有限制
	// 1: sortFields=field1&sortFields=field2&sortFields=field3
	// 2: sortFields=field1,field2,field3

	sortFields := []string{}
	raw, ok := r.URL.Query()["sortFields"] // 先获取
	if !ok {
		return filter
	}

	// 循环判断是单个还是多个
	for _, v := range raw {
		if v == "" {
			continue
		}
		// 如果包含逗号，存在多个做拆分
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					sortFields = append(sortFields, p)
				}
			}
		} else {
			// 单个值
			sortFields = append(sortFields, v)
		}
	}

	filter.SortFields = sortFields

	return filter
}

// ReadJSON 读取body参数绑定到dst上
func (app *AppToolkit) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	err := gotk.ReadJSON(w, r, dst)
	if err != nil {
		a := errs.ErrBadRequest.WithMessage(err.Error())
		fmt.Println(a.Error())
		app.FAIL(w, r, errs.ErrBadRequest.WithMessage(err.Error()))
		return false
	}

	slog.InfoContext(r.Context(), "arg", slog.Any("arg", dst))
	return true
}

// ReadJSONAndCheck 读取body参数绑定到dst上并校验，dst是实现了gotk.Verifiyer接口类型的指针
func (app *AppToolkit) ReadJSONAndCheck(w http.ResponseWriter, r *http.Request, dst gotk.Verifiyer) bool {
	// 解码绑定
	ok := app.ReadJSON(w, r, dst)
	if !ok {
		return false
	}

	// 校验
	if v, ok := gotk.DoVerifiy(dst); ok {
		return true
	} else {
		msg := v.GetOne()
		err := errs.ErrBadRequest.WithMessage(msg)
		app.FAIL(w, r, err)
		return false
	}
}

// FAIL 写入请求失败的方法，data数据返回null
func (app *AppToolkit) FAIL(w http.ResponseWriter, r *http.Request, a *gotk.ApiError) {
	if a == nil {
		slog.DebugContext(r.Context(), "AppToolkit.FAIL 请提供 *gotk.ApiError")
		a = errs.ErrServerError
	}
	app.write(w, r, a, nil)
}

// SUCC 写入请求成功的方法
func (app *AppToolkit) SUCC(w http.ResponseWriter, r *http.Request, data any) {
	app.write(w, r, errs.ErrOK, data)
}

// write 通用写入的方法
func (app *AppToolkit) write(
	w http.ResponseWriter,
	r *http.Request,
	a *gotk.ApiError,
	data any,
) {
	if a == nil {
		slog.DebugContext(r.Context(), "AppToolkit.write 请提供 *gotk.ApiError")
		a = errs.ErrOK
	}

	var url = r.Method + " - " + r.RequestURI
	var status = "请求失败"

	var logwrite = slog.ErrorContext
	if http.StatusOK == a.StatusCode() {
		logwrite = slog.InfoContext
		status = "请求成功"
	}
	// TODO: 通过上下文获取用户id
	logwrite(
		r.Context(),
		status,
		"url",
		url,
		"msg",
		slog.StringValue(a.Error()),
		"data",
		slog.AnyValue(data),
	)

	// 写入响应
	err := gotk.WriteJSON(w, r, a, data)

	// 写入失败
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			"写入响应失败: "+url,
			"msg",
			slog.StringValue(a.Error()),
			"data",
			slog.AnyValue(data),
			slog.String("write_error", err.Error()),
		)
	}
}
