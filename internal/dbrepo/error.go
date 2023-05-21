package dbrepo

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNotFound          = errors.New("记录不存在")
	ErrISBNAlreadyExists = errors.New("isbn 已经存在")
	ErrUpdateFailed      = errors.New("更新失败")
	ErrDelFailed         = errors.New("删除失败")
)

func IsCustomDBError(err error) bool {
	switch {
	case errors.Is(err, ErrNotFound),
		errors.Is(err, ErrISBNAlreadyExists),
		errors.Is(err, ErrUpdateFailed),
		errors.Is(err, ErrDelFailed):
		return true
	}
	return false
}

// 统一处理mysql错误
func handleMySQLError(err error) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		var mysqlErr *mysql.MySQLError
		// 重复键错误, 具体哪个字段重复，由具体业务判断
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// NOTE: unq_isbn 是数据库定义唯一索引名字
			if strings.Contains(err.Error(), "unq_isbn") {
				return ErrISBNAlreadyExists
			}
			// NOTE: 根据业务添加其他的
		}
	}

	return err
}
