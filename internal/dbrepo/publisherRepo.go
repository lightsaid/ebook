package dbrepo

import (
	"github.com/lightsaid/ebook/internal/models"
)

type PublisherRepo interface {
	Create(name string) (uint64, error)
	Update(name string) error
	Get(id uint64) (*models.Publisher, error)
	List() ([]*models.Publisher, error)
	Delete(id uint64) error
}

var _ PublisherRepo = (*publisherRepo)(nil)

type publisherRepo struct {
	DB Queryable
}

func NewPublisherRepo(db Queryable) *publisherRepo {
	var repo = &publisherRepo{
		DB: db,
	}
	return repo
}

func (r *publisherRepo) Create(name string) (uint64, error) {
	sql := `insert publisher set publisher_name = ?;`
	result, err := r.DB.Exec(sql, name)
	return insertErrorHandler(result, err)
}

func (r *publisherRepo) Update(name string) error {
	sql := `update publisher set publisher_name = ? where deleted_at is not null;`
	result, err := r.DB.Exec(sql, name)
	return updateErrorHandler(result, err)
}

func (r *publisherRepo) Get(id uint64) (publisher *models.Publisher, err error) {
	sql := `
		select 
			id, publisher_name, created_at, updated_at 
		from 
			publisher 
		where 
			id = ? and deleted_at is null;`
	publisher = new(models.Publisher)
	err = r.DB.Get(publisher, sql, id)
	return publisher, err
}

func (r *publisherRepo) List() (list []*models.Publisher, err error) {
	sql := `select * from publisher where deleted_at is null;`
	err = r.DB.Select(&list, sql)
	return list, err
}

func (r *publisherRepo) Delete(id uint64) error {
	sql := `update publisher set deleted_at = now() where id = ?;`
	result, err := r.DB.Exec(sql, id)
	return updateErrorHandler(result, err)
}
