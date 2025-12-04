package dbrepo

import (
	"context"
	"log/slog"

	"github.com/lightsaid/ebook/internal/models"
)

type PublisherRepo interface {
	Create(ctx context.Context, name string) (uint64, error)
	Update(ctx context.Context, name string) error
	Get(ctx context.Context, id uint64) (*models.Publisher, error)
	List(ctx context.Context) ([]*models.Publisher, error)
	Delete(ctx context.Context, id uint64) error
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

func (r *publisherRepo) Create(ctx context.Context, name string) (uint64, error) {
	sql := `insert publisher set publisher_name = ?;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query, slog.String("name", name))

	result, err := r.DB.ExecContext(ctx, query, name)

	return insertErrorHandler(result, err)
}

func (r *publisherRepo) Update(ctx context.Context, name string) error {
	sql := `update publisher set publisher_name = ? where deleted_at is null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query, slog.String("name", name))

	result, err := r.DB.ExecContext(ctx, query, name)

	return updateErrorHandler(result, err)
}

func (r *publisherRepo) Get(ctx context.Context, id uint64) (publisher *models.Publisher, err error) {
	sql := `
		select 
			id, publisher_name, created_at, updated_at 
		from 
			publisher 
		where 
			id = ? and deleted_at is null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, spaceRex.ReplaceAllString(query, " "), slog.Int64("id", int64(id)))

	publisher = new(models.Publisher)

	err = r.DB.GetContext(ctx, publisher, query, id)
	return publisher, err
}

func (r *publisherRepo) List(ctx context.Context) (list []*models.Publisher, err error) {
	sql := `select * from publisher where deleted_at is null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query)

	err = r.DB.SelectContext(ctx, &list, query)
	return list, err
}

func (r *publisherRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update publisher set deleted_at = now() where id = ?;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query, slog.Int64("id", int64(id)))

	result, err := r.DB.ExecContext(ctx, query, id)
	return updateErrorHandler(result, err)
}
