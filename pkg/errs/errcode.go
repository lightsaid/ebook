package errs

import (
	"net/http"

	"github.com/lightsaid/gotk"
)

var (
	ErrOK                  = gotk.NewApiError(http.StatusOK, "10000", "请求成功")
	ErrBadRequest          = gotk.NewApiError(http.StatusBadRequest, "10400", "入参错误")
	ErrUnauthorized        = gotk.NewApiError(http.StatusUnauthorized, "10401", "请先登陆")
	ErrForbidden           = gotk.NewApiError(http.StatusForbidden, "10403", "禁止访问")
	ErrNotFound            = gotk.NewApiError(http.StatusNotFound, "10404", "查无此数据")
	ErrMethodNotAllowed    = gotk.NewApiError(http.StatusMethodNotAllowed, "10405", "请求方法不支持")
	ErrRecordExists        = gotk.NewApiError(http.StatusConflict, "10409", "数据已存在")
	ErrUnprocessableEntity = gotk.NewApiError(http.StatusUnprocessableEntity, "10422", "请求无法处理")
	ErrTooManyRequests     = gotk.NewApiError(http.StatusTooManyRequests, "10429", "请求繁忙")
	ErrServerError         = gotk.NewApiError(http.StatusInternalServerError, "10500", "请求错误，请稍后重试！")
)
