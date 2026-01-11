package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lightsaid/ebook/internal/config"
	"github.com/lightsaid/ebook/internal/dbcache"
	"github.com/lightsaid/ebook/internal/dbrepo"
)

func (app *Application) serve(logger *slog.Logger) error {
	srv := http.Server{
		Addr:        fmt.Sprintf("0.0.0.0:%d", app.config.ServerPort),
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

		log.Println("执行释放资源操作")
		dbrepo.Close()
		dbcache.Close()

		shutdownError <- nil
	}()

	fmt.Println("start api server on ", srv.Addr)
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

// Healthcheck 服务健康检查
func (app *Application) Healthcheck(w http.ResponseWriter, r *http.Request) {
	app.SUCC(w, r, "请求成功")
}

// ReloadConfig 重新配置加载配置
//
// 例如某段时间异常错误非常多，修改配置，重新加载配置文件，使用新的配置方便调试
// 如 改变log级别等
// 虽然可以重载配置，并不是所有配置都会生效
// 比如连接的数据库参数，已经连接了，那就不会在使用，除非从新建立连接
func (app *Application) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: 设计一下鉴权
	config.Load(&app.config, app.envFiles...)
	confbuf, _ := json.MarshalIndent(&app.config, "", " ")
	fmt.Println(string(confbuf))
	app.SUCC(w, r, string(confbuf))
}
