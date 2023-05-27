package dbrepo

import (
	"context"
	"database/sql"
)

type Repository struct {
	Book BookRepo
}

// Queryable 提取 sql.DB 和 sql.Tx 公共的方法当作一个接口
type Queryable interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// NOTE: 采用这这种方式，是为了在执行事务上可以复用代码，更加灵活性
// NewRepository 创建一个Repository仓库，使用 Queryable 接口，同时兼容 sql.DB 和 sql.Tx 接口
func NewRepository(db Queryable) *Repository {
	return &Repository{
		Book: NewBookRepo(db),
	}
}

// func NewRepository(db *sql.DB) Repository {
// 	return Repository{
// 		Book: NewBookRepo(db),
// 	}
// }
