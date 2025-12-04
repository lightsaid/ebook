package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

// PostBookHandler godoc
//
//	@Summary		创建图书
//	@Description	用户创建图书
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.Book	true	"EBook payload"
//	@Success		200		{object}	int
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/book [post]
func (app *Application) PostBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if ok := app.ShouldBindJSONAndCheck(w, r, &book); !ok {
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
	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	var book models.Book
	if ok := app.ShouldBindJSONAndCheck(w, r, &book); !ok {
		return
	}

	book.ID = id

	// TODO: 获取/赋值才更新

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
