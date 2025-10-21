package dbrepo

import (
	"github.com/lightsaid/ebook/internal/models"
)

type BookCategoryRepo interface {
	Create(bc models.BookCategory) (uint64, error)
	ListByBookID(bookID uint64) (list []*models.BookCategory, err error)
	ListByCategoryID(categoryID uint64) (list []*models.BookCategory, err error)
	DeleteByBookID(bookID uint64) error
	DeleteByCategoryID(bookID uint64) error

	BatchInsert(list []models.BookCategory) error
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

func (r *bookCategoryRepo) Create(bc models.BookCategory) (uint64, error) {
	sql := `insert book_categories set book_id=?, category_id=?;`
	result, err := r.DB.Exec(sql, bc.BookID, bc.CategoryID)
	return insertErrorHandler(result, err)
}

func (r *bookCategoryRepo) BatchInsert(list []models.BookCategory) error {
	sql := `insert into book_categories(book_id, category_id) values `
	args := []any{}
	for i, x := range list {
		if i > 0 {
			sql += ","
		}
		sql += "(?, ?)"
		args = append(args, x.BookID, x.CategoryID)
	}

	ctx, cancel := makeCtx()
	defer cancel()

	_, err := insertErrorHandler(r.DB.ExecContext(ctx, sql, args...))

	return err
}

func (r *bookCategoryRepo) ListByBookID(bookID uint64) (list []*models.BookCategory, err error) {
	sql := `select * from book_categories where book_id=?;`
	err = r.DB.Select(&list, sql, bookID)
	return list, err
}

func (r *bookCategoryRepo) ListByCategoryID(categoryID uint64) (list []*models.BookCategory, err error) {
	sql := `select * from book_categories where category_id=?;`
	err = r.DB.Select(&list, sql, categoryID)
	return list, err
}

func (r *bookCategoryRepo) DeleteByBookID(bookID uint64) error {
	sql := `delete from book_categories where book_id=?`
	result, err := r.DB.Exec(sql, bookID)
	return updateErrorHandler(result, err)
}

func (r *bookCategoryRepo) DeleteByCategoryID(bookID uint64) error {
	sql := `delete from book_categories where category_id=?`
	result, err := r.DB.Exec(sql, bookID)
	return updateErrorHandler(result, err)
}
