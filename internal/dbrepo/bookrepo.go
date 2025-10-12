package dbrepo

import (
	"context"

	"github.com/lightsaid/ebook/internal/models"
)

type BookRepo interface {
	Create(book *models.Book) (uint64, error)
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
		r.BookRepo.Create(nil)
		return nil
	})
	return nil
}
