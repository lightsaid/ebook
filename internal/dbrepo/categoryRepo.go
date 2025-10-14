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
	sql := `insert category set category_name=?, icon=?, sort=?;`
	result, err := r.DB.Exec(sql, category.CategoryName, category.Icon, category.Sort)
	return insertErrorHandler(result, err)
}

func (r *categoryRepo) Update(category models.Category) error {
	sql := `update category set category_name=?, icon=?, sort=? where id = ? and deleted_at is null;`
	result, err := r.DB.Exec(
		sql,
		category.CategoryName,
		category.Icon,
		category.Sort,
		category.ID,
	)
	return updateErrorHandler(result, err)
}

func (r *categoryRepo) Get(id uint64) (*models.Category, error) {
	sql := `
		select
			id, category_name, icon, sort, created_at, updated_at 
		from 
			category 
		where 
			id = ? and deleted_at is null;`

	category := new(models.Category)
	err := r.DB.Get(category, sql, id)
	return category, err
}

func (r *categoryRepo) List() (list []*models.Category, err error) {
	sql := `select * from category where deleted_at is null;`
	err = r.DB.Select(&list, sql)
	return list, err
}

func (r *categoryRepo) Delete(id uint64) error {
	sql := `update category set deleted_at = now() where id = ?;`
	result, err := r.DB.Exec(sql, id)
	return updateErrorHandler(result, err)
}
