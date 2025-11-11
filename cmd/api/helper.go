package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Application) readIDParam(r *http.Request, field string) (int64, error) {

	id, err := strconv.ParseInt(chi.URLParam(r, field), 10, 64)
	if err != nil || id < 1 {
		// TODO: 错误处理
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
