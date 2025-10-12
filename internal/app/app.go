package app

// import (
// 	"database/sql"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/lightsaid/ebook/internal/config"
// 	"github.com/lightsaid/ebook/internal/dbrepo"
// 	"github.com/lightsaid/ebook/internal/fileupload"
// 	"github.com/lightsaid/ebook/pkg/logger"
// )

// // application 一个管理API服务的结构体
// type application struct {
// 	cfg      config.AppConfig
// 	store    *dbrepo.Repository
// 	uploader fileupload.FileUploader
// }

// func Serve() error {
// 	return newApplication().serve()
// }

// func newApplication() *application {
// 	// 解析命令行参数
// 	var cfgPath string
// 	flag.StringVar(&cfgPath, "path", "configs/app.json", "config 配置文件路径")
// 	flag.Parse()

// 	// 加载app配置
// 	cfg, err := config.LoadAppConfig(cfgPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// 格式化输出一下配置，仅开发模式下
// 	cfg.Println()

// 	// 设置全局日志
// 	logger.SetGlobalLogger()

// 	// 连接mysql
// 	db, err := sql.Open("mysql", cfg.DSN)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// ping 确定能连接上数据库
// 	if err = db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 设置最大链接数、最大空闲数和最大空闲时间
// 	db.SetMaxOpenConns(cfg.MaxOpenConns)
// 	db.SetMaxIdleConns(cfg.MaxIdleConns)
// 	db.SetConnMaxIdleTime(cfg.MaxIdleTimeToDuration())

// 	// 文件上传接口
// 	uploader := fileupload.NewLocalUplader(cfg.UploadPath, cfg.AllowsExt, cfg.MaxUploadByte)

// 	a := &application{
// 		cfg:      cfg,
// 		store:    dbrepo.NewRepository(db),
// 		uploader: uploader,
// 	}

// 	return a
// }

// func (a *application) serve() error {
// 	var address = fmt.Sprintf("0.0.0.0:%d", a.cfg.Port)

// 	// 创建 http.Server 用来启动 Web Server
// 	srv := http.Server{
// 		Addr:         address,
// 		Handler:      a.routes(),
// 		IdleTimeout:  time.Minute,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 30 * time.Second,
// 	}

// 	logger.InfoLog.Println("Starting API Server on: " + address)
// 	return srv.ListenAndServe()
// }
