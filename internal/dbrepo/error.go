package dbrepo

import (
	"database/sql"
	"errors"
)

var ErrInsertIDZero = errors.New("新增数据返回的id为0")

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
	eff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if eff <= 0 {
		return sql.ErrNoRows
	}
	return nil
}
