package dbrepo

import "github.com/lightsaid/ebook/internal/models"

type CategoryRepo interface {
	Create(category models.Category) (uint64, error)
	Update(category models.Category) error
	Get(id uint64) (*models.Category, error)
	List() ([]*models.Category, error)
	Delete(id uint64) error
}

var _ CategoryRepo = (*categoryRepo)(nil)

type categoryRepo struct {
	DB Queryable
}

func NewCategoryRepo(db Queryable) *categoryRepo {
	var repo = &categoryRepo{
		DB: db,
	}
	return repo
}

func (r *categoryRepo) Create(category models.Category) (uint64, error) {
	panic("TODO:")
}
func (r *categoryRepo) Update(category models.Category) error {
	panic("TODO:")
}
func (r *categoryRepo) Get(id uint64) (*models.Category, error) {
	panic("TODO:")

}
func (r *categoryRepo) List() ([]*models.Category, error) {
	panic("TODO:")
}
func (r *categoryRepo) Delete(id uint64) error {
	panic("TODO:")
}
