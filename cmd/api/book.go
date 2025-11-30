package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

// PostBookHandler
func (app *Application) PostBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if ok := app.ShouldBindJSON(w, r, &book); !ok {
		return
	}

	newID, err := app.Db.BookRepo.CreateTx(r.Context(), &book)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, newID)
}

// GetBookHandler
func (app *Application) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	book, err := app.Db.BookRepo.Get(r.Context(), id)
	if err != nil {
		a = dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, book)
}

// PutBookHandler
func (app *Application) PutBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if ok := app.ShouldBindJSON(w, r, &book); !ok {
		return
	}

	err := app.Db.BookRepo.UpdateTx(r.Context(), &book)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

}

// DeleteBookHandler
func (app *Application) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := app.Db.BookRepo.Delete(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// ListBookHandler
func (app *Application) ListBookHandler(w http.ResponseWriter, r *http.Request) {
	filter := app.readPageQuery(r)
	dataVo, err := app.Db.BookRepo.ListWithCategory(r.Context(), filter)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, dataVo)
}
