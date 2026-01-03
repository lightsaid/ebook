package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/ebook/internal/config"
	"github.com/lightsaid/ebook/internal/dbrepo"
)

var tDb *sqlx.DB
var tRepo dbrepo.Repository

// 执行test先执行
func TestMain(m *testing.M) {
	// 输出一下执行文件的路径，如果下面配置不对，可以通过os.Args从新拼接
	fmt.Println(os.Args)

	var conf config.DbConfig
	config.Load(&conf, "../../../configs/develop.env")
	var err error
	tDb, err = dbrepo.Open(conf)
	if err != nil {
		log.Fatal("cannot open db:", err)
	}

	tRepo = dbrepo.NewRepository(tDb)

	// 执行测试
	code := m.Run()

	// 关闭资源
	dbrepo.Close()

	os.Exit(code)
}
