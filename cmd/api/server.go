package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lightsaid/ebook/internal/dbrepo"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	// Chi中间件执行顺序=洋葱模型（Onion Model）
	// Use() 注册的顺序 = 执行的顺序（从外往里）
	// 返回时 = 从里往外（反方向）

	router.Use(middleware.CleanPath)

	router.Get("/v1/healthcheck", app.Healthcheck)
	router.Get("/v1/books", app.ListBookHandler)

	mux := chi.NewRouter()
	mux.Mount("/api", router)

	return mux
	// 超时控制
	//	return http.TimeoutHandler(mux, 5*time.Second, "请求超时")
}

func (app *Application) serve(logger *slog.Logger) error {
	srv := http.Server{
		Addr:        "0.0.0.0:6000",
		Handler:     app.routes(),
		IdleTimeout: time.Minute,

		// 从客户端读取请求头和body超时设置
		ReadTimeout: 10 * time.Second,

		// WriteTimeout 从请求头读完时就已经开始计时；
		// 对于大多数请求，当进入处理请求函数时，就开始计时了；
		// 因此如果执行时间超过WriteTimeout，还没写入响应，链接就会被关闭，但是处理函数不会被终止，请求方，没有任何返回。
		WriteTimeout: 15 * time.Second,

		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		// syscall.SIGINT 中断信号。通常由用户在终端按下 Ctrl+C 时发送给前台进程。
		// syscall.SIGTERM 终止信号。这是标准的、通用的程序终止请求信号，不指定信号类型的 kill 命令默认发送此信号。
		// Kubernetes、systemd 等现代服务管理系统在停止服务时，通常也会先发送 SIGTERM。
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		log.Println("接受到系统信号: ", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("优雅关机")

		// 关机将在不中断任何活动连接的情况下优雅地关闭服务器
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// TODO: 释放资源

		log.Println("执行释放资源操作")
		dbrepo.Close()

		shutdownError <- nil
	}()

	log.Println("start api server on ", srv.Addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) Healthcheck(w http.ResponseWriter, r *http.Request) {
	// 测试 优雅关机
	slog.DebugContext(r.Context(), ">>> start healthcheck")
	// time.Sleep(2 * time.Second)
	slog.DebugContext(r.Context(), ">>> write response")
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(envelope{"status": "ok"})
	if err != nil {
		log.Println(r.RequestURI, " write response fail: ", err)
	}
}

func (app *Application) ListBookHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info(r.URL.Path)
	list, err := app.Db.BookRepo.ListWithCategory(10, 30)
	if err != nil {
		slog.Info(err.Error())
		// TODO:
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
