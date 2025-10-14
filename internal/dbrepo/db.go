package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Queryable提取sql.DB和sql.Tx公共的方法当作一个接口,
// 为了执行事务可以调用sql.DB的操作，在事务中服用基础的CRUD方法。
type Queryable interface {
	// sql
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// sqlx
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

var (
	Db *sqlx.DB
)

const (
	defaultDuration = 3 * time.Second
)

func Open() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root.cc",
		"127.0.0.1",
		3306,
		"db_ebook",
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(15)

	Db = db

	return db, nil
}

func Close() {
	log.Println("closeing db")
	if err := Db.Close(); err != nil {
		log.Println("Close DB: ", err)
	}
}

func execTx(ctx context.Context, query Queryable, fn func(Repository) error) error {
	db, ok := query.(*sqlx.DB)
	if !ok {
		return errors.New("`query Queryable` is not sqlx.DB")
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	repository := NewRepository(tx)
	err = fn(repository)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w rb err: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultDuration)
}
