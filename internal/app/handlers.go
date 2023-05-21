package app

import (
	"net/http"
	"strconv"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/logger"
)

func (a *application) addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.methodNotAllowed(w, http.MethodPost)
		return
	}
	var book = new(models.Book)
	err := a.readJSON(w, r, book)
	if err != nil {
		a.writeJSON(w, http.StatusBadRequest, msgWrapp(err.Error()))
		return
	}

	// TODO: 验证参数

	id, err := a.store.Book.Insert(book)
	if err != nil {
		logger.ErrorfoLog.Println("addBook failed -> ", err)
		a.errorResponse(w, err)
		return
	}
	a.writeJSON(w, http.StatusOK, wrapper{"id": id})
}

func (a *application) getBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.methodNotAllowed(w, http.MethodGet)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		a.writeJSON(w, http.StatusBadRequest, msgWrapp("id 无效"))
		return
	}
	book, err := a.store.Book.Get(uint(id))
	if err != nil {
		logger.ErrorfoLog.Println("getBook failed -> ", err.Error(), "id = ", id)
		a.errorResponse(w, err)
		return
	}

	a.writeJSON(w, http.StatusOK, dataWrapp(book))
}

func (a *application) putBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		a.methodNotAllowed(w, http.MethodPut)
		return
	}
	var book models.Book
	if err := a.readJSON(w, r, &book); err != nil {
		a.writeJSON(w, http.StatusBadRequest, msgWrapp(err.Error()))
		return
	}
	err := a.store.Book.Update(book.ID, &book)
	if err != nil {
		logger.ErrorfoLog.Println("putbook failed -> ", err)
		a.errorResponse(w, err)
		return
	}
	a.writeJSON(w, http.StatusOK, msgWrapp("更新成功"))
}

func (a *application) delBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		a.methodNotAllowed(w, http.MethodDelete)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		a.writeJSON(w, http.StatusBadRequest, msgWrapp("id 无效"))
		return
	}

	err = a.store.Book.Del(uint(id))
	if err != nil {
		logger.ErrorfoLog.Println("delBook failed -> ", err, "id = ", id)
		a.errorResponse(w, err)
		return
	}
	a.writeJSON(w, http.StatusOK, msgWrapp("删除成功"))
}

func (a *application) listBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.methodNotAllowed(w, http.MethodGet)
		return
	}
	// 获取查询参数 list?page=1&size=10
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	p := dbrepo.NewPager(page, size)
	books, err := a.store.Book.List(p)
	if err != nil {
		logger.ErrorfoLog.Println("listBook failed -> ", err, "page = ", p)
		a.errorResponse(w, err)
		return
	}
	a.writeJSON(w, http.StatusOK, books)
}

func (a *application) bulkInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.methodNotAllowed(w, http.MethodPost)
		return
	}
	var books []*models.Book
	if err := a.readJSON(w, r, &books); err != nil {
		a.errorResponse(w, err)
		return
	}

	logger.InfoLog.Println("insert books num -> ", len(books))
	if len(books) == 0 {
		a.writeJSON(w, http.StatusBadRequest, msgWrapp("参数不能为空"))
		return
	}

	aff, err := a.store.Book.BulkInsert(books)
	if err != nil {
		a.errorResponse(w, err)
		return
	}
	total := len(books)
	failed := total - int(aff)
	a.writeJSON(w, http.StatusOK, wrapper{"total": total, "ok": aff, "failed": failed})
}

func (a *application) testTx(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			a.writeJSON(w, http.StatusBadRequest, msgWrapp("id 无效"))
			return
		}
		err = a.store.Book.TestTx(uint(id))
		if err != nil {
			a.errorResponse(w, err)
			return
		}
		a.writeJSON(w, http.StatusOK, msgWrapp("操作成功"))
		return
	}
	a.methodNotAllowed(w, http.MethodPost, http.MethodGet)
}
