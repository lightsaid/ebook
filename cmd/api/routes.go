package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lightsaid/gotk"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	// Chi中间件执行顺序=洋葱模型（Onion Model）
	// Use() 注册的顺序 = 执行的顺序（从外往里）
	// 返回时 = 从里往外（反方向）

	router.Use(middleware.CleanPath)

	router.Get("/v1/healthcheck", app.Healthcheck)

	{
		// 图书api
		router.Post("/v1/book", app.PostBookHandler)
		router.Get("/v1/book/{id:[0-9]+}", app.GetBookHandler)
		router.Put("/v1/book/{id:[0-9]+}", app.PutBookHandler)
		router.Delete("/v1/book/{id:[0-9]+}", app.DeleteBookHandler)
		router.Get("/v1/books", app.ListBookHandler)
	}

	{
		// 作者api
		router.Post("/v1/author", app.PostAuthorHandler)
		router.Get("/v1/author/{id:[0-9]+}", app.GetAuthorHandler)
		router.Put("/v1/author/{id:[0-9]+}", app.PutAuthorHandler)
		router.Delete("/v1/author/{id:[0-9]+}", app.DeleteAuthorHandler)
		router.Get("/v1/authors", app.ListAuthorHandler)
	}

	{
		// 分类api
		router.Post("/v1/category", app.PostCategoryHandler)
		router.Get("/v1/category/{id:[0-9]+}", app.GetCategoryHandler)
		router.Put("/v1/category/{id:[0-9]+}", app.PutCategoryHandler)
		router.Delete("/v1/category/{id:[0-9]+}", app.DeleteCategoryHandler)
		router.Get("/v1/categories", app.ListCategoryHandler)
	}

	{
		// 出版社api
		router.Post("/v1/publisher", app.PostPublisherHandler)
		router.Get("/v1/publisher/{id:[0-9]+}", app.GetPublisherHandler)
		router.Put("/v1/publisher/{id:[0-9]+}", app.PutPublisherHandler)
		router.Delete("/v1/publisher/{id:[0-9]+}", app.DeletePublisherHandler)
		router.Get("/v1/publishers", app.ListPublisherHandler)
	}

	mux := chi.NewRouter()
	mux.Mount("/api", router)

	return gotk.SetRequestIDCtx(mux)
	// 超时控制
	//	return http.TimeoutHandler(mux, 5*time.Second, "请求超时")
}
