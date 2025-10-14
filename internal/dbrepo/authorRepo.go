package dbrepo

import (
	"github.com/lightsaid/ebook/internal/models"
)

type AuthorRepo interface {
	Create(authorName string) (uint64, error)
	Update(id uint64, authorName string) error
	Get(id uint64) (*models.Author, error)
	List() ([]*models.Author, error)
	Delete(id uint64) error
}

var _ AuthorRepo = (*authorRepo)(nil)

type authorRepo struct {
	DB Queryable
}

func NewAuthorRepo(db Queryable) *authorRepo {
	var repo = &authorRepo{
		DB: db,
	}
	return repo
}

func (r *authorRepo) Create(authorName string) (uint64, error) {
	sql := `insert author set author_name = ?;`
	result, err := r.DB.Exec(sql, authorName)
	return insertErrorHandler(result, err)
}

func (r *authorRepo) Update(id uint64, authorName string) error {
	sql := `update author set author_name = ? where id =? and deleted_at is null;`
	result, err := r.DB.Exec(sql, authorName, id)
	return updateErrorHandler(result, err)
}

func (r *authorRepo) Get(id uint64) (author *models.Author, err error) {
	sql := `
		select 
			id, author_name, created_at, updated_at 
		from 
			author 
		where 
			id = ? and deleted_at is null;`
	author = new(models.Author)
	err = r.DB.Get(author, sql, id)
	return author, err
}

func (r *authorRepo) List() (list []*models.Author, err error) {
	sql := `select * from author where deleted_at is null;`
	err = r.DB.Select(&list, sql)
	return list, err
}

func (r *authorRepo) Delete(id uint64) error {
	sql := `update author set deleted_at = now() where id = ?;`
	result, err := r.DB.Exec(sql, id)
	return updateErrorHandler(result, err)
}
