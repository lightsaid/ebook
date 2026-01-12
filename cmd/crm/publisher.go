package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

// PostPublisherHandler godoc
//
//	@Summary		添加出版社
//	@Description	添加一个出版社到管理系统
//	@Tags			Publisher
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		 models.Publisher	true	"添加出版社请求体"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/publisher [post]
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

// GetPublisherHandler godoc
//
// @Summary 获取一个出版社
// @Description 根据id获取出版社
//
//	@Tags			Publisher
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"publisher id"
//	@Success		200	{object}	ApiResponse{data=models.Publisher}
//	@Router			/v1/publisher/{id} [get]
func (app *Application) GetPublisherHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}
	app.SUCC(w, r, id)
}

// PutPublisherHandler godoc
//
//	@Summary		更新出版社
//	@Description	根据id更新出版社
//	@Tags			Publisher
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.Publisher	true	"更新出版社请求体"
//	@Param			id		path		int				true	"publisher id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/publisher/{id} [put]
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

// DeletePublisherHandler godoc
//
//	@Summary		删除一个出版社
//	@Description	根据id删除一个出版社
//	@Tags			Publisher
//	@Produce		json
//	@Param			id		path		int				true	"publisher id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/publisher/{id} [delete]
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

// ListPublisherHandler godoc
//
//	@Summary		获取出版列表
//	@Description	获取出版社列表，没分页
//	@Tags			Publisher
//	@Produce		json
//	@Success		200			{object}	ApiResponse{data=[]models.Publisher}
//	@Router			/v1/publishers [get]
func (app *Application) ListPublisherHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.PublisherRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
