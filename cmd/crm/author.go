package main

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/gotk"
)

type AuthorNameReq struct {
	AuthorName string `json:"authorName"`
}

func (p AuthorNameReq) Verifiy(v *gotk.Validator) {
	v.Check(p.AuthorName != "", "authorName", "作者名称称必填")
}

// PostAuthorHandler godoc
//
//	@Summary		添加作者
//	@Description	添加一个作者到管理系统
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			author	body		AuthorNameReq	true	"添加作者请求体"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/author [post]
func (app *Application) PostAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var author AuthorNameReq
	if ok := app.ReadJSONAndCheck(w, r, &author); !ok {
		return
	}

	newID, err := store.AuthorRepo.Create(r.Context(), author.AuthorName)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, newID)
}

// GetAuthorHandler godoc
//
//	@Summary		获取作者
//	@Description	根据id获取一个作者
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"作者id"
//	@Success		200	{object}	ApiResponse{data=models.Author}
//	@Router			/v1/author/{id} [get]
func (app *Application) GetAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	author, err := store.AuthorRepo.Get(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, author)
}

// PutAuthorHandler godoc
//
//	@Summary		更新作者
//	@Description	添加一个作者到管理系统
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			author	body		AuthorNameReq	true	"添加作者请求体"
//	@Param			id		path		int				true	"作者id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/author/{id} [put]
func (app *Application) PutAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var author models.Author
	if ok := app.ReadJSONAndCheck(w, r, &author); !ok {
		return
	}

	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := store.AuthorRepo.Update(r.Context(), id, author.AuthorName)
	if err != nil {
		a = dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// DeleteAuthorHandler godoc
//
//	@Summary		删除作者
//	@Description	根据id删除一个作者
//	@Tags			Author
//	@Produce		json
//	@Param			id		path		int				true	"作者id"
//	@Success		200		{object}	ApiResponse{data=int}
//	@Router			/v1/author/{id} [delete]
func (app *Application) DeleteAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, a := app.ReadIntParam(r, "id")
	if a != nil {
		app.FAIL(w, r, a)
		return
	}

	err := store.AuthorRepo.Delete(r.Context(), id)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, id)
}

// 分页查询，比如： /api/v1/authors?pageNum=1&pageSize=5&sortFields=-created_at&sortFields=id
//
// 注意 sortFields 可排序字段别写错，错误示例：sortFields=-createdAt，要以数据库字段为准，
// 这里不对字段做转换
// GetAuthorHandler godoc
//
//	@Summary		获取作者列表
//	@Description	分页获取作者列表
//	@Tags			Author
//	@Produce		json
//	@Param			pageNum		query		int			false	"页码"
//	@Param			pageSize	query		int			false	"每页多少条"
//	@Param			sortFields	query		[]string	false	"排序字段"
//	@Success		200			{object}	ApiResponse{data=dbrepo.PageQueryVo{list=[]models.Author}}
//	@Router			/v1/authors [get]
func (app *Application) ListAuthorHandler(w http.ResponseWriter, r *http.Request) {
	filter := app.ReadPageQuery(r)

	data, err := store.AuthorRepo.List(r.Context(), filter)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, data)
}
