package dbrepo

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
)

var (
	ErrInsertIDZero       = errors.New("新增数据返回的id为0")
	ErrBookCategoryNoRows = errors.New("图书分类不存在")

	ErrNotFound     = errors.New("记录不存在")
	ErrInsertFailed = errors.New("插入数据失败")
	ErrNoEffectDB   = errors.New("qb not is *sql.DB")
	ErrNoEmail      = errors.New("邮箱地址不能为空")
)

// ConvertToApiError 将db错误转换为 *gotk.ApiError
func ConvertToApiError(err error) *gotk.ApiError {
	// MySQL 错误码参照：
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html

	if errors.Is(err, ErrNotFound) {
		return errs.ErrNotFound.WithError(err)
	}
	if errors.Is(err, ErrInsertFailed) || errors.Is(err, ErrNoEffectDB) {
		return errs.ErrServerError.WithError(err)
	}

	if err == sql.ErrNoRows {
		return errs.ErrNotFound.WithError(err)
	}

	dbErr, ok := err.(*mysql.MySQLError)
	if ok {
		if dbErr.Number == 1062 { // 数据冲突，已存在
			return errs.ErrRecordExists.WithError(err)
		}
	}

	return errs.ErrServerError.WithError(err)
}
