package app

import (
	"net/http"

	"github.com/lightsaid/ebook/internal/config"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/book", a.getBook)         // GET    获取一个图书 book?id=1
	mux.HandleFunc("/v1/book/bulk", a.bulkInsert) // POST 批量添加
	mux.HandleFunc("/v1/book/post", a.addBook)    // POST   创建一个图书
	mux.HandleFunc("/v1/book/put", a.putBook)     // PUT    更新一个图书
	mux.HandleFunc("/v1/book/del", a.delBook)     // DELETE 删除一个图书 del?id=1
	mux.HandleFunc("/v1/book/list", a.listBook)   // GET    获取图书列表 list?page=1&size=10

	// 一些测试路由
	if a.cfg.Env == config.Env_Dev {
		mux.HandleFunc("/v1/test/tx", a.testTx)    // GET/POST tx?id=1 测试事物路由
		mux.HandleFunc("/v1/upload", a.uploadFile) // GET/POST 测试文件上传接口

	}

	// 定义一个文件服务
	fileServer := http.FileServer(http.Dir("templates/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
