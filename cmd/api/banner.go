package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

// PostBannerHandler
func (app *Application) PostBannerHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Banner
	if ok := app.ShouldBindJSON(w, r, &p); !ok {
		return
	}

	id, err := app.Db.BannerRepo.Create(r.Context(), p.BannerName)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// GetBannerHandler
func (app *Application) GetBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}
	app.SUCC(w, r, id)
}

// PutBannerHandler
func (app *Application) PutBannerHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Banner
	if ok := app.ShouldBindJSON(w, r, &p); !ok {
		return
	}

	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := app.Db.BannerRepo.Update(r.Context(), p.BannerName)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// DeleteBannerHandler
func (app *Application) DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.readIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := app.Db.BannerRepo.Delete(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)

}

// ListBannerHandler
func (app *Application) ListBannerHandler(w http.ResponseWriter, r *http.Request) {
	list, err := app.Db.BannerRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
