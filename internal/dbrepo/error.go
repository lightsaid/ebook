package dbrepo

import (
	"database/sql"
	"errors"
	"log"

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
)

// insertErrorHandler 公共处理新增数据的错误
func insertErrorHandler(result sql.Result, err error) (uint64, error) {
	if err != nil {
		return 0, err
	}
	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if newID <= 0 {
		return 0, ErrInsertIDZero
	}
	return uint64(newID), nil
}

// updateErrorHandler 公共处理更新数据的错误
func updateErrorHandler(result sql.Result, err error) error {
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Println("updateErrorHandler RowsAffected fail: ", err)
		return err
	}
	return nil

	// eff, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// // NOTE: 如果数据存在，但是并没有改变任何值，result.RowsAffected()也是会返回0的
	// // 通常在更新之前都会先查询一遍，因此这里可以允许为0
	// if eff <= 0 {
	// 	return sql.ErrNoRows
	// }
	// return nil
}

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
