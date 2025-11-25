package main

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

func (app *Application) readIntParam(r *http.Request, field string) (uint64, *gotk.ApiError) {
	val := chi.URLParam(r, field)
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil || id < 1 {
		errVal := fmt.Errorf("无效的 %s 参数: %s", field, val)
		return 0, errs.ErrBadRequest.With(errVal, errVal.Error())
	}
	return uint64(id), nil
}

func (app *Application) readPageQuery(r *http.Request) dbrepo.Filters {
	var filter dbrepo.Filters
	pageNumText := r.URL.Query().Get("pageNum")
	pageSizeText := r.URL.Query().Get("pageSize")

	// 无需处理错误，0值使用时，会设置为默认值
	filter.PageNum, _ = strconv.Atoi(pageNumText)
	filter.PageSize, _ = strconv.Atoi(pageSizeText)

	// 需要支持两种传参数方式：
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

func (app *Application) ShouldBindJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	err := gotk.ReadJSON(w, r, dst)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.WithMessage(err.Error()))
		return false
	}

	slog.InfoContext(r.Context(), "arg", slog.Any("arg", dst))
	return true
}

// ShouldBindJSONAndCheck绑定请求入参并执行校验
func (app *Application) ShouldBindJSONAndCheck(w http.ResponseWriter, r *http.Request, dst gotk.Verifiyer) bool {
	// 解码绑定
	ok := app.ShouldBindJSON(w, r, dst)
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

// FAIL 请求失败
func (app *Application) FAIL(w http.ResponseWriter, r *http.Request, a *gotk.ApiError) {
	if a == nil {
		slog.InfoContext(r.Context(), "由于(a *gotk.ApiError)没值,设为ErrServerError")
		a = errs.ErrServerError
	}
	app.write(w, r, a, nil)
}

// SUCC 请求成功
func (app *Application) SUCC(w http.ResponseWriter, r *http.Request, data any) {
	app.write(w, r, errs.ErrOK, data)
}

func (app *Application) write(
	w http.ResponseWriter,
	r *http.Request,
	a *gotk.ApiError,
	data any,
) {
	if a == nil {
		slog.InfoContext(r.Context(), "由于(a *gotk.ApiError)没值,设为ErrOK")
		a = errs.ErrOK
	}

	var message = r.Method + " - " + r.RequestURI

	// 如果要写入的状态是200，成功的日志
	if http.StatusOK == a.StatusCode() {
		slog.InfoContext(
			r.Context(),
			message,
			"status",
			slog.IntValue(a.StatusCode()),
			"bizCode",
			slog.StringValue(a.BizCode()),
		)
	} else {
		// 记录错误日志
		slog.ErrorContext(
			r.Context(),
			message,
			"status",
			slog.IntValue(a.StatusCode()),
			"bizCode",
			slog.StringValue(a.BizCode()),
			"data",
			slog.AnyValue(data),
		)
	}

	// 写入响应
	err := gotk.WriteJSON(w, r, a, data)

	// 写入失败
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			"写入响应失败: "+message,
			"status",
			slog.IntValue(a.StatusCode()),
			"bizCode",
			slog.StringValue(a.BizCode()),
			"data",
			slog.AnyValue(data),
			slog.String("api_error", a.Error()),
			slog.String("write_error", err.Error()),
		)
	}
}
