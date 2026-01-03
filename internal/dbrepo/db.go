package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/ebook/internal/config"
)

var (
	Db       *sqlx.DB
	spaceRex = regexp.MustCompile(`\s+`)

	// 内部使用的工具,提取出来，方便使用，不用记忆多个工具函数名字
	// 通过 dbtk.xxx 智能提示即可
	dbtk = toolkit{}
)

type envelop map[string]any

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

func Open(conf config.DbConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DbUser,
		conf.DbPswd,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
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

// baseRepo 基础接口定义在这里
type baseRepo interface {
	// defaultSortSafelist 导出默认的安全排序字段
	defaultSortSafelist() []string
}

// 内部使用的工具箱 toolkit
type toolkit struct{}

// debugSQL 使用sqlx.Named 和 db.Rebind 处理sql并输出sql日志，返回新的sql和参数
// 因此sql语句中的参数必须是sqlx :param 格式参数
//
// arg 必须为结构体或者map(key为数据库字段)，struct(带db tag)
func (*toolkit) debugSQL(ctx context.Context, db Queryable, sql string, arg any) (string, []any, error) {
	// 使用 sqlx.Named 把 :param 统一转换为'?'占位符，并生成 args slice
	query, args, err := sqlx.Named(sql, arg)
	if err != nil {
		slog.DebugContext(ctx, "debugSQL->sqlx.Named(sql,arg)", "error", err)
		return "", nil, err
	}

	// 将'?'根据数据驱动绑定占位符类型(?/$1/@p1...)
	query = db.Rebind(query)

	// 清理多余的\t、\n,输出更容易阅读
	cleanSQL := spaceRex.ReplaceAllString(query, " ")

	// 输出 log
	slog.DebugContext(ctx, "debugSQL SQL", slog.String("sql", cleanSQL))
	slog.DebugContext(ctx, "debugSQL Args", slog.Any("args", args))

	return cleanSQL, args, nil
}

// execTx 执行事务公共方法
func (*toolkit) execTx(ctx context.Context, query Queryable, fn func(Repository) error) error {
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

// withTimeout 设置超时，并返回新的context
func (*toolkit) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}

// insertErrorHandler 公共处理新增数据的错误
func (*toolkit) insertErrorHandler(ctx context.Context, result sql.Result, err error) (uint64, error) {
	if err != nil {
		slog.DebugContext(ctx, "dbtk.updateErrorHandler: "+err.Error())
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
func (*toolkit) updateErrorHandler(ctx context.Context, result sql.Result, err error) error {
	if err != nil {
		slog.DebugContext(ctx, "dbtk.updateErrorHandler: "+err.Error())
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		slog.DebugContext(ctx, "dbtk.updateErrorHandler RowsAffected fail: "+err.Error())
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

// dbtk.calculateMetadata 计算分页数据
func (*toolkit) calculateMetadata(totalCount, pageNum, pageSize int) Metadata {
	if totalCount == 0 {
		return Metadata{}
	}

	return Metadata{
		PageNum:    pageNum,
		PageSize:   pageSize,
		LastPage:   (totalCount + pageSize - 1) / pageSize,
		TotalCount: totalCount,
	}
}

// dbtk.makePageQueryVo 构建统一返回数据,extras 是额外数据，有则返回第一个
func (*toolkit) makePageQueryVo(metadata Metadata, list any, extras ...any) *PageQueryVo {
	var extra any
	if len(extras) > 0 {
		extra = extras[0]
	}
	return &PageQueryVo{
		List:      list,
		Metadata:  metadata,
		ExtraData: extra,
	}
}
