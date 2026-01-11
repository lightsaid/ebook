package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/errs"
)

// PostCategoryHandler 处理创建分类
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

// GetCategoryHandler 获取一个分类
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

// PutCategoryHandler 更新分类
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

// DeleteCategoryHandler 删除一个分类
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

func (app *Application) ListCategoryHandler(w http.ResponseWriter, r *http.Request) {
	list, err := store.CategoryRepo.List(r.Context())
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, list)
}
