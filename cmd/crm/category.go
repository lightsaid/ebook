package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/errs"
)

// PostCategoryHandler godoc
//
//	@Summary		添加分类
//	@Description	添加一个分类到管理系统
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		 models.Category	true	"添加分类请求体"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/category [post]
func (app *Application) PostCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category

	// 绑定并检验入参
	ok := app.ReadJSONAndCheck(w, r, &c)
	if !ok {
		return
	}

	// 执行分类创建
	newID, err := store.CategoryRepo.Create(r.Context(), c)
	if err != nil {
		// 处理数据库错误
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	// 创建成功响应
	app.SUCC(w, r, newID)
}

// GetCategoryHandler godoc
//
// @Summary 获取一个分类
// @Description 根据id获取分类
//
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"category id"
//	@Success		200	{object}	ApiResponse{data=models.Category}
//	@Router			/v1/banner/{id} [get]
func (app *Application) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadIntParam(r, "id")
	if err != nil {
		app.FAIL(w, r, err)
		return
	}

	category, err2 := store.CategoryRepo.Get(r.Context(), id)
	if err2 != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, category)
}

// PutCategoryHandler godoc
//
//	@Summary		更新分类
//	@Description	根据id更新分类
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.Category	true	"更新分类请求体"
//	@Param			id		path		int				true	"category id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/category/{id} [put]
func (app *Application) PutCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, aerr := app.ReadIntParam(r, "id")
	if aerr != nil {
		app.FAIL(w, r, aerr)
		return
	}

	var c models.Category

	// 绑定参数并校验
	ok := app.ReadJSONAndCheck(w, r, &c)
	if !ok {
		return
	}

	// 根据id获取分类
	category, err := store.CategoryRepo.Get(r.Context(), id)
	if err != nil {
		aerr := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, aerr)
		return
	}

	// 赋值
	category.CategoryName = c.CategoryName
	category.Sort = c.Sort
	if c.Icon != "" {
		category.Icon = c.Icon
	}

	// 更新
	err = store.CategoryRepo.Update(r.Context(), *category)
	if err != nil {
		// 处理数据库错误
		aerr := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, aerr)
		return
	}
	app.SUCC(w, r, "更新成功")

}

// DeleteCategoryHandler godoc
//
//	@Summary		删除一个分类
//	@Description	根据id删除一个分类
//	@Tags			Category
//	@Produce		json
//	@Param			id		path		int				true	"category id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/category/{id} [delete]
func (app *Application) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, aerr := app.ReadIntParam(r, "id")
	if aerr != nil {
		app.FAIL(w, r, aerr)
		return
	}

	// 先获取，在删除
	_, err := store.CategoryRepo.Get(r.Context(), id)
	if err != nil {
		aerr = dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, aerr)
		return
	}

	// 执行删除
	err = store.CategoryRepo.Delete(r.Context(), id)
	if err != nil {
		aerr = dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, aerr)
		return
	}

	app.SUCC(w, r, errs.ErrOK)
}

// ListCategoryHandler godoc
//
//	@Summary		获取分类列表
//	@Description	获取分类列表，没分页
//	@Tags			Category
//	@Produce		json
//	@Success		200			{object}	ApiResponse{data=[]models.Category}
//	@Router			/v1/categories [get]
func (app *Application) ListCategoryHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.CategoryRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
