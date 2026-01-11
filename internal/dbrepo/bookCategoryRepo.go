package dbrepo

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lightsaid/ebook/internal/models"
)

type BookCategoryRepo interface {
	Create(ctx context.Context, bc models.BookCategory) (uint64, error)
	ListByBookID(ctx context.Context, bookID uint64) (list []*models.BookCategory, err error)
	ListByCategoryID(ctx context.Context, categoryID uint64) (list []*models.BookCategory, err error)
	DeleteByBookID(ctx context.Context, bookID uint64) error
	DeleteByCategoryID(ctx context.Context, bookID uint64) error

	BatchInsert(ctx context.Context, list []models.BookCategory) error
}

var _ BookCategoryRepo = (*bookCategoryRepo)(nil)

type bookCategoryRepo struct {
	DB Queryable
}

func NewBookCategoryRepo(db Queryable) *bookCategoryRepo {
	var repo = &bookCategoryRepo{
		DB: db,
	}
	return repo
}

func (r *bookCategoryRepo) Create(ctx context.Context, bc models.BookCategory) (uint64, error) {
	sql := r.DB.Rebind(`insert book_categories set book_id=?, category_id=?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	slog.DebugContext(
		ctx, sql,
		"book_id",
		slog.Int64Value(int64(bc.BookID)),
		"category_id", slog.Int64Value(int64(bc.CategoryID)),
	)

	result, err := r.DB.ExecContext(ctx, sql, bc.BookID, bc.CategoryID)
	return dbtk.insertErrorHandler(ctx, result, err)
}

func (r *bookCategoryRepo) BatchInsert(ctx context.Context, list []models.BookCategory) error {
	if len(list) <= 0 {
		return ErrBookCategoryNoRows
	}
	sql := `insert into book_categories(book_id, category_id) values `
	parts := make([]string, len(list))
	args := []any{}
	for i, x := range list {
		parts[i] = "(?, ?)"
		args = append(args, x.BookID, x.CategoryID)
	}
	sql += strings.Join(parts, ",")

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)
	slog.InfoContext(ctx, spaceRex.ReplaceAllString(query, " "), "args", slog.AnyValue(args))

	// NOTE: 这里只会返回影响行数，因此不能使用 dbtk.insertErrorHandler
	// _, err := dbtk.insertErrorHandler(r.DB.ExecContext(ctx, sql, args...))
	result, err := r.DB.ExecContext(ctx, query, args...)
	return dbtk.updateErrorHandler(ctx, result, err)
}

func (r *bookCategoryRepo) ListByBookID(ctx context.Context, bookID uint64) (list []*models.BookCategory, err error) {
	sql := r.DB.Rebind(`select * from book_categories where book_id=?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	slog.DebugContext(ctx, sql, "book_id", slog.Int64Value(int64(bookID)))

	err = r.DB.SelectContext(ctx, &list, sql, bookID)
	return list, err
}

func (r *bookCategoryRepo) ListByCategoryID(ctx context.Context, categoryID uint64) (list []*models.BookCategory, err error) {
	sql := r.DB.Rebind(`select * from book_categories where category_id=?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	slog.DebugContext(ctx, sql, "category_id", slog.Int64Value(int64(categoryID)))

	err = r.DB.SelectContext(ctx, &list, sql, categoryID)
	return list, err
}

func (r *bookCategoryRepo) DeleteByBookID(ctx context.Context, bookID uint64) error {
	sql := r.DB.Rebind(`delete from book_categories where book_id=?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	slog.DebugContext(ctx, sql, "book_id", slog.Int64Value(int64(bookID)))

	result, err := r.DB.ExecContext(ctx, sql, bookID)
	return dbtk.updateErrorHandler(ctx, result, err)
}

func (r *bookCategoryRepo) DeleteByCategoryID(ctx context.Context, categoryID uint64) error {
	sql := r.DB.Rebind(`delete from book_categories where category_id=?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	slog.DebugContext(ctx, sql, "category_id", slog.Int64Value(int64(categoryID)))

	result, err := r.DB.ExecContext(ctx, sql, categoryID)
	return dbtk.updateErrorHandler(ctx, result, err)
}
