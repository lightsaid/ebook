package dbrepo

import (
	"context"
	"log/slog"

	"github.com/lightsaid/ebook/internal/models"
)

type CategoryRepo interface {
	Create(ctx context.Context, category models.Category) (uint64, error)
	Update(ctx context.Context, category models.Category) error
	Get(ctx context.Context, id uint64) (*models.Category, error)
	List(ctx context.Context) ([]*models.Category, error)
	Delete(ctx context.Context, id uint64) error
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

func (r *categoryRepo) Create(ctx context.Context, category models.Category) (uint64, error) {
	// NOTE: 这是mysql特有的语法？好像是，记得，并且sqlx.Named不支持
	// sql := `insert category set category_name=?, icon=?, sort=?;`
	sql := `insert into category(category_name, icon, sort) values (:category_name, :icon, :sort);`

	ctx, cancal := dbtk.withTimeout(ctx)
	defer cancal()

	query, args, err := dbtk.debugSQL(ctx, r.DB, sql, category)
	if err != nil {
		return 0, err
	}
	result, err := r.DB.ExecContext(ctx, query, args...)

	return insertErrorHandler(result, err)
}

func (r *categoryRepo) Update(ctx context.Context, category models.Category) error {
	sql := `update category set category_name=:category_name, icon=:icon, sort=:sort where id=:id and deleted_at is null;`

	ctx, cancal := dbtk.withTimeout(ctx)
	defer cancal()

	query, args, err := dbtk.debugSQL(ctx, r.DB, sql, category)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(
		ctx,
		query,
		args...,
	)
	return updateErrorHandler(result, err)
}

func (r *categoryRepo) Get(ctx context.Context, id uint64) (*models.Category, error) {
	sql := `
		select
			id, category_name, icon, sort, created_at, updated_at 
		from 
			category 
		where 
			id=? and deleted_at is null;`

	ctx, cancal := dbtk.withTimeout(ctx)
	defer cancal()

	sql = r.DB.Rebind(sql)

	cleanSQL := spaceRex.ReplaceAllString(sql, " ")

	slog.InfoContext(ctx, cleanSQL, "id", id)

	category := new(models.Category)
	err := r.DB.GetContext(ctx, category, sql, id)
	return category, err
}

func (r *categoryRepo) List(ctx context.Context) (list []*models.Category, err error) {
	sql := `select * from category where deleted_at is null;`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	sql = r.DB.Rebind(sql)

	slog.InfoContext(ctx, sql)

	err = r.DB.SelectContext(ctx, &list, sql)
	return list, err
}

func (r *categoryRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update category set deleted_at = now() where id = ?;`

	sql = r.DB.Rebind(sql)

	slog.InfoContext(ctx, sql, "id", id)

	result, err := r.DB.Exec(sql, id)

	return updateErrorHandler(result, err)
}
