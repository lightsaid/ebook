package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

// PostPublisherHandler
func (app *Application) PostPublisherHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Publisher
	if ok := app.ReadJSONAndCheck(w, r, &p); !ok {
		return
	}

	id, err := store.PublisherRepo.Create(r.Context(), p.PublisherName)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// GetPublisherHandler
func (app *Application) GetPublisherHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}
	app.SUCC(w, r, id)
}

// PutPublisherHandler
func (app *Application) PutPublisherHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Publisher
	if ok := app.ReadJSONAndCheck(w, r, &p); !ok {
		return
	}

	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	// TODO: 获取/赋值才更新

	err := store.PublisherRepo.Update(r.Context(), p.PublisherName)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// DeletePublisherHandler
func (app *Application) DeletePublisherHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := store.PublisherRepo.Delete(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)

}

// ListPublisherHandler
func (app *Application) ListPublisherHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.PublisherRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
