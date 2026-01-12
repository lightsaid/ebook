package main

import (
	"net/http"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/internal/types"
)

// PostBannerHandler godoc
//
//	@Summary		添加banner
//	@Description	添加一个banner到管理系统
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param			banner	body		 models.Banner	true	"添加banner请求体"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/banner [post]
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

// GetBannerHandler godoc
//
// @Summary 获取banner
// @Description 根据id获取banner
//
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"banner id"
//	@Success		200	{object}	ApiResponse{data=models.Banner}
//	@Router			/v1/banner/{id} [get]
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

// PutBannerHandler godoc
//
//	@Summary		更新banner
//	@Description	根据id更新banner
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param			banner	body		models.Banner	true	"更新banner请求体"
//	@Param			id		path		int				true	"banner id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/banner/{id} [put]
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

// DeleteBannerHandler godoc
//
//	@Summary		删除banner
//	@Description	根据id删除一个banner
//	@Tags			Banner
//	@Produce		json
//	@Param			id		path		int				true	"banner id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/banner/{id} [delete]
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

// ListBannerHandler godoc
//
//	@Summary		获取banner列表
//	@Description	获取banner列表，没分页
//	@Tags			Banner
//	@Produce		json
//	@Success		200			{object}	ApiResponse{data=[]models.Banner}
//	@Router			/v1/banners [get]
func (app *Application) ListBannerHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.BannerRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
