package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lightsaid/gotk"

	docs "github.com/lightsaid/ebook/docs/crm"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	// Chi中间件执行顺序=洋葱模型（Onion Model）
	// Use() 注册的顺序 = 执行的顺序（从外往里）
	// 返回时 = 从里往外（反方向）

	router.Use(middleware.CleanPath)
	router.Use(gotk.WithRequestIDCtx)
	router.Use(app.EnableCORS)
	router.Use(app.AccessLog)
	router.Use(app.RecoverPanic)

	router.Get("/v1/healthcheck", app.Healthcheck)
	router.Post("/v1/signin", app.SignIn)
	router.Post("/v1/reset_pswd", app.RestPassword)
	router.Post("/v1/renew_token", app.RenewAccessToken)

	router.Group(func(r chi.Router) {
		r.Use(app.RequiredAuth)

		{
			r.Post("/v1/reload/config", app.ReloadConfig)
		}

		{ // user api
			r.Post("/v1/profile", app.UpdateProfile)

			r.Get("/v1/users", app.GetListUser)
		}

		{ // book api

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

		{
			// banner api
			router.Post("/v1/banner", app.PostBannerHandler)
			router.Get("/v1/banner/{id:[0-9]+}", app.GetBannerHandler)
			router.Put("/v1/banner/{id:[0-9]+}", app.PutBannerHandler)
			router.Delete("/v1/banner/{id:[0-9]+}", app.DeleteBannerHandler)
			router.Get("/v1/banners", app.ListBannerHandler)
		}

		{ // shoppingCart api

		}

		{ // order api

		}
	})

	mux := chi.NewRouter()
	mux.Mount("/api", router)

	app.setupSwaggerDoc(mux)

	return mux
	// 超时控制
	//	return http.TimeoutHandler(mux, 5*time.Second, "请求超时")
}

func (app *Application) setupSwaggerDoc(mux *chi.Mux) {
	fmt.Printf("Swagger: http://localhost:%d/swagger/index.html\n", app.config.ServerPort)

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", app.config.ServerPort) // swagger服务host
	docs.SwaggerInfo.BasePath = "/api"                                         // api请求前缀

	docsURL := fmt.Sprintf("http://%s/swagger/doc.json", docs.SwaggerInfo.Host)

	// fmt.Println("docsURL: ", docsURL)

	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(docsURL),
	))

	// 输出注册的路由
	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s -> %s\n", method, route)
		return nil
	})
}
