package dbrepo

import (
	"context"
	"log/slog"

	"github.com/lightsaid/ebook/internal/models"
)

type BannerRepo interface {
	Create(ctx context.Context, banner *models.Banner) (uint64, error)
	Update(ctx context.Context, banner *models.Banner) error
	Get(ctx context.Context, id uint64) (*models.Publisher, error)
	List(ctx context.Context) ([]*models.Publisher, error)
	Delete(ctx context.Context, id uint64) error
}

var _ BannerRepo = (*bannerRepo)(nil)

type bannerRepo struct {
	DB Queryable
}

func NewBannerRepo(db Queryable) *bannerRepo {
	var repo = &bannerRepo{
		DB: db,
	}
	return repo
}

func (r *bannerRepo) Create(ctx context.Context, banner *models.Banner) (uint64, error) {
	query := `insert into banners(slogan,link_type,link_url,image_url,enable,sort)
		values(:slogan, :link_type, :link_url, :image_url, :enable, :sort)`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query, args, err := debugSQL(ctx, r.DB, query, banner)
	if err != nil {
		return 0, err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)

	return insertErrorHandler(result, err)
}

func (r *bannerRepo) Update(ctx context.Context, banner *models.Banner) error {
	query := `update banners set 
	slogan=:slogan, link_type=:link_type, 
	link_url=:link_url, image_url=:image_url, 
	enable=:enable, sort=:sort
	where deleted_at is not null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query, args, err := debugSQL(ctx, r.DB, query, banner)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)

	return updateErrorHandler(result, err)
}

func (r *bannerRepo) Get(ctx context.Context, id uint64) (banners *models.Publisher, err error) {
	sql := `select * from banners where id = ? and deleted_at is null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	banners = new(models.Publisher)

	err = r.DB.GetContext(ctx, banners, query, id)
	return banners, err
}

func (r *bannerRepo) List(ctx context.Context) (list []*models.Publisher, err error) {
	sql := `select * from banners where deleted_at is null;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query)

	err = r.DB.SelectContext(ctx, &list, query)
	return list, err
}

func (r *bannerRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update banners set deleted_at = now() where id = ?;`

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query, slog.Int64("id", int64(id)))

	result, err := r.DB.ExecContext(ctx, query, id)
	return updateErrorHandler(result, err)
}
