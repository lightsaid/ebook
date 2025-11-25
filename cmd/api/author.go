package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lightsaid/ebook/internal/dbrepo"
)

func (app *Application) PostAuthorHandler(w http.ResponseWriter, r *http.Request)   {}
func (app *Application) GetAuthorHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) PutAuthorHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) DeleteAuthorHandler(w http.ResponseWriter, r *http.Request) {}

// 分页查询，比如： /api/v1/authors?pageNum=1&pageSize=5&sortFields=-created_at&sortFields=id
//
// 注意 sortFields 可排序字段别写错，错误示例：sortFields=-createdAt，要以数据库字段为准，
// 这里不对字段做转换
func (app *Application) ListAuthorHandler(w http.ResponseWriter, r *http.Request) {
	filter := app.readPageQuery(r)

	data, err := app.Db.AuthorRepo.List(r.Context(), filter)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	by, _ := json.Marshal(data)
	fmt.Println(">>>> ", string(by))

	app.SUCC(w, r, data)
}
