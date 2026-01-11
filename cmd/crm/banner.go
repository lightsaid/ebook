package main

import (
	"net/http"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/internal/types"
)

// PostBannerHandler
func (app *Application) PostBannerHandler(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	if ok := app.ReadJSONAndCheck(w, r, &banner); !ok {
		return
	}

	id, err := store.BannerRepo.Create(r.Context(), &banner)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// GetBannerHandler
func (app *Application) GetBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}
	data, err := store.BannerRepo.Get(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}
	app.SUCC(w, r, data)
}

// PutBannerHandler
func (app *Application) PutBannerHandler(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	if ok := app.ReadJSONAndCheck(w, r, &banner); !ok {
		return
	}

	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	b, err := store.BannerRepo.Get(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	b.Slogan = banner.Slogan
	b.LinkType = banner.LinkType
	b.LinkUrl = banner.LinkUrl
	b.ImageUrl = banner.ImageUrl
	b.Enable = banner.Enable
	b.Sort = banner.Sort
	b.UpdatedAt = types.GxTime{Time: time.Now()}

	err = store.BannerRepo.Update(r.Context(), b)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// DeleteBannerHandler
func (app *Application) DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := store.BannerRepo.Delete(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)

}

// ListBannerHandler
func (app *Application) ListBannerHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.BannerRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
