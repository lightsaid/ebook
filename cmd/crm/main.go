package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/lightsaid/ebook/internal/config"
	"github.com/lightsaid/ebook/internal/dbcache"
	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/ebook/pkg/apptk"
	"github.com/lightsaid/ebook/pkg/logger"
	"github.com/lightsaid/gotk"
)

var (
	store dbrepo.Repository
	cache dbcache.Repository
)

type Application struct {
	apptk.AppToolkit
	jwt      gotk.TokenMaker
	envFiles types.ArrayString
	config   struct {
		config.CRMConfig
		config.DbConfig
		config.JWTConfig
		config.RedisConfig
	}
}

func main() {

	// 解析命令行参数获取配置文件
	var envFiles types.ArrayString
	flag.Var(&envFiles, "env", "配置文件，支持指定多个")
	flag.Parse()

	app := Application{}

	// 保存配置文件路径，可用户刷新
	app.envFiles = envFiles

	// 解析配置数据到app.config
	config.Load(&app.config, envFiles...)

	confbuf, _ := json.MarshalIndent(&app.config, "", " ")
	fmt.Println(string(confbuf))

	// 初始化日志输出实例
	instance := logger.NewLogger(os.Stdout, app.config.LogLevel, gotk.TextType)
	slog.SetDefault(instance)

	jwt, err := gotk.NewJWTMaker(
		app.config.JWTConfig.SecretKey,
		app.config.JWTConfig.Issuer,
	)

	if err != nil {
		log.Fatalln(err)
	}

	app.jwt = jwt

	// 与数据库建立连接
	conn, err := dbrepo.Open(app.config.DbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 创建数据crud实例
	store = dbrepo.NewRepository(conn)

	// 与redis建立连接
	rdb, err := dbcache.Open(app.config.RedisConfig)
	if err != nil {
		log.Fatalln(err)
	}

	cache = dbcache.NewRepository(rdb)

	// 启动接口服务
	if err := app.serve(instance); err != nil {
		log.Fatalln(err)
	}
}
