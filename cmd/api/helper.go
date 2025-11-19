package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
)

func (app *Application) readIntParam(r *http.Request, field string) (int64, *gotk.ApiError) {
	val := chi.URLParam(r, field)
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil || id < 1 {
		errVal := fmt.Errorf("无效的 %s 参数: %s", field, val)
		return 0, errs.ErrBadRequest.With(errVal, errVal.Error())
	}
	return id, nil
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
func (app *Application) FAIL(w http.ResponseWriter, r *http.Request, err *gotk.ApiError) {
	werr := gotk.WriteJSON(w, r, err, nil)
	slog.ErrorContext(
		r.Context(),
		r.Method+"-"+r.RequestURI,
		slog.String("err", err.Error()),
		"write_err",
		werr,
	)
}

// SUCC 请求成功 TODO:
func (app *Application) SUCC(w http.ResponseWriter, r *http.Request, data any) {
	_ = gotk.WriteJSON(w, r, errs.ErrOK, data)
}
