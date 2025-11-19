package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

func (app *Application) PostCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	ok := app.ShouldBindJSONAndCheck(w, r, &c)
	if !ok {
		return
	}
	newID, err := app.Db.CategoryRepo.Create(c)
	if err != nil {
		// 处理数据库错误
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}
	app.SUCC(w, r, newID)
}

func (app *Application) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.FAIL(w, r, err)
		return
	}

	category, err2 := app.Db.CategoryRepo.Get(uint64(id))
	if err2 != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, category)
}

func (app *Application) PutCategoryHandler(w http.ResponseWriter, r *http.Request) {}

func (app *Application) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {}

func (app *Application) ListCategoryHandler(w http.ResponseWriter, r *http.Request) {
	list, err := app.Db.CategoryRepo.List()
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
