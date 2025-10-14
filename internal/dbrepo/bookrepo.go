package dbrepo

import (
	"context"

	"github.com/lightsaid/ebook/internal/models"
)

type BookRepo interface {
	Create(book *models.Book) (uint64, error)
	CreateTx(book *models.Book) (uint64, error)
	Get(id uint64) (*models.Book, error)
	Update(book *models.Book) error   // 仅更新books表
	UpdateTx(book *models.Book) error // 更新图书和与之关联的分类、出版社、作者
	List() ([]*models.Book, error)
	ListByCategory(categoryID uint64) ([]*models.Book, error)
	Delete(id uint64) error
}

var _ BookRepo = (*bookRepo)(nil)

type bookRepo struct {
	DB Queryable
}

func NewBookRepo(db Queryable) *bookRepo {
	var repo = &bookRepo{
		DB: db,
	}
	return repo
}

func (r *bookRepo) Create(book *models.Book) (uint64, error) {
	sql := `
	insert into books set 
		isbn=:isbn, 
		title=:title, 
		subtitle=:subtitle,
		author_id=:author_id,
		cover_url=:cover_url, 
		publisher_id=:publisher_id, 
		pubdate=:pubdate,
		price=:price, 
		status=:status, 
		type=:type,
		stock=:stock,
		source_url=:source_url,
		description=:description
	`

	args := map[string]any{
		"isbn":         book.ISBN,
		"title":        book.Title,
		"subtitle":     book.Subtitle,
		"author_id":    book.AuthorID,
		"cover_url":    book.CoverUrl,
		"publisher_id": book.PublisherID,
		"pubdate":      book.Pubdate,
		"price":        book.Price,
		"status":       book.Status,
		"type":         book.Type,
		"stock":        book.Stock,
		"source_url":   book.SourceUrl,
		"description":  book.Description,
	}

	ctx, cancel := makeCtx()
	defer cancel()

	result, err := r.DB.NamedExecContext(ctx, sql, args)
	return insertErrorHandler(result, err)
}

func (r *bookRepo) CreateTx(book *models.Book) (uint64, error) {
	return 0, nil
}

func (repo *bookRepo) Get(id uint64) (book *models.Book, err error) {
	return nil, nil
}

func (repo *bookRepo) Update(book *models.Book) error {
	return nil
}

func (repo *bookRepo) UpdateTx(book *models.Book) error {
	return nil
}

func (repo *bookRepo) List() ([]*models.Book, error) {
	return nil, nil
}

func (repo *bookRepo) ListByCategory(categoryID uint64) ([]*models.Book, error) {
	return nil, nil
}

func (repo *bookRepo) Delete(id uint64) error {
	execTx(context.Background(), repo.DB, func(r Repository) error {
		// r.BookRepo.Create(nil)
		return nil
	})
	return nil
}
