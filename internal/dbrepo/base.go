package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// dbBase 内部实现通用方法，嵌套在repo，提供给repo使用
type dbBase struct{}

// execTx 定义个执行事务公共的方法
func (dbBase) execTx(ctx context.Context, qb Queryable, fn func(*Repository) error) error {
	// qb => sql.DB/sql.Tx
	db, ok := qb.(*sql.DB)
	if !ok {
		return errors.New("b.DB not is *sql.DB")
	}
	// 开启事务
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	r := NewRepository(tx)
	if _, err = r.Book.Get(6); err != nil {
		return err
	}

	q := NewRepository(tx)
	if err = fn(q); err != nil {
		// 执行不成功，则 Rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 执行成功 Commit
	return tx.Commit()
}
