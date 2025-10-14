package tests

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/ebook/internal/dbrepo"
)

var tDb *sqlx.DB
var tRepo dbrepo.Repository

// 执行test先执行
func TestMain(m *testing.M) {
	var err error
	tDb, err = dbrepo.Open()
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
