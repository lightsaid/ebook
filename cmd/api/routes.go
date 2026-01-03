package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lightsaid/gotk"

	docs "github.com/lightsaid/ebook/docs/api"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	{
		// banner api
		router.Post("/v1/banner", app.PostBannerHandler)
		router.Get("/v1/banner/{id:[0-9]+}", app.GetBannerHandler)
		router.Put("/v1/banner/{id:[0-9]+}", app.PutBannerHandler)
		router.Delete("/v1/banner/{id:[0-9]+}", app.DeleteBannerHandler)
		router.Get("/v1/banners", app.ListBannerHandler)
	}

	{
		// 用户api
		router.Post("/v1/user/register", app.UserRegisterHandler)
		router.Post("/v1/user/login", app.UserLoginHandler)
		router.Put("/v1/user/update", app.UserUpdateHandler)
		router.Get("/v1/user/profile", app.GetUserProfile)
		router.Post("/v1/user/renewToken", app.RenewTokenHandler)
	}

	{
		//  订单api
		router.Post("/v1/order", app.PostOrderHandler)
		router.Post("/v1/order/pay", app.PayOrderHandler)
		router.Get("/v1/order/{id:[0-9]+}", app.GetOrderHandler)
		router.Put("/v1/order/{id:[0-9]+}", app.PutOrderHandler)
		router.Delete("/v1/order/{id:[0-9]+}", app.DeleteOrderHandler)
		router.Get("/v1/orders", app.ListOrderHandler)
	}

	{
		// 购物车api
		router.Post("/v1/shopping/cart", app.PostShoppingCartHandler)
		router.Delete("/v1/shopping/cart/{id:[0-9]+}", app.DeleteShoppingCartHandler)
		router.Get("/v1/shopping/carts", app.ListShoppingCartHandler)
	}

	mux := chi.NewRouter()
	mux.Mount("/api", router)

	setupSwaggerDoc(mux)

	return gotk.WithRequestIDCtx(mux)
	// 超时控制
	//	return http.TimeoutHandler(mux, 5*time.Second, "请求超时")
}

// TODO: 配置

func setupSwaggerDoc(mux *chi.Mux) {
	slog.Info("Swagger: http://localhost:4000/swagger/index.html")

	docs.SwaggerInfo.Host = "localhost:4000" // swagger服务host
	docs.SwaggerInfo.BasePath = "/api"       // api请求前缀

	docsURL := fmt.Sprintf("http://%s/swagger/doc.json", docs.SwaggerInfo.Host)

	// fmt.Println("docsURL: ", docsURL)

	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(docsURL),
	))

	// 调试输出注册的路由
	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s -> %s\n", method, route)
		return nil
	})
}
