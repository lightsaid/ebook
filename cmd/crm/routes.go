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

// TODO: 登录和使用中间件查询用户信息，使用redis缓存用户信息5分钟，
// 下一次请求，直接冲redis里获取用户信息，如果没有再从mysql查询并存储到redis里。

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	// Chi中间件执行顺序=洋葱模型（Onion Model）
	// Use() 注册的顺序 = 执行的顺序（从外往里）
	// 返回时 = 从里往外（反方向）

	router.Use(middleware.CleanPath)

	router.Get("/v1/healthcheck", app.Healthcheck)
	router.Post("/v1/signin", app.SignIn)
	router.Post("/v1/reset_pswd", app.RestPassword)

	router.Group(func(r chi.Router) {
		r.Use(app.RequiredAuth)

		{
			r.Post("/v1/reload/config", app.ReloadConfig)
		}

		{ // user api
			r.Post("/v1/renew_token", app.RenewAccessToken)
			r.Post("/v1/profile", app.UpdateProfile)

			r.Get("/v1/users", app.GetListUser)
		}

		{ // book api

		}

		{ // author api

		}

		{ // category api

		}

		{ // publisher api

		}

		{ // banner api

		}

		{ // shoppingCart api

		}

		{ // order api

		}
	})

	mux := chi.NewRouter()
	mux.Mount("/api", router)

	app.setupSwaggerDoc(mux)

	return gotk.WithRequestIDCtx(mux)
	// 超时控制
	//	return http.TimeoutHandler(mux, 5*time.Second, "请求超时")
}

func (app *Application) setupSwaggerDoc(mux *chi.Mux) {
	fmt.Printf("Swagger: http://localhost:%d/swagger/index.html", app.config.ServerPort)

	docs.SwaggerInfo.Host = fmt.Sprintf("0.0.0.0:%d", app.config.ServerPort) // swagger服务host
	docs.SwaggerInfo.BasePath = "/api"                                       // api请求前缀

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
