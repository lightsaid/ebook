package dbrepo

import (
	"context"
	"log/slog"

	"github.com/lightsaid/ebook/internal/models"
)

type BannerRepo interface {
	Create(ctx context.Context, banner *models.Banner) (uint64, error)
	Update(ctx context.Context, banner *models.Banner) error
	Get(ctx context.Context, id uint64) (*models.Banner, error)
	List(ctx context.Context) ([]*models.Banner, error)
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

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query, args, err := dbtk.debugSQL(ctx, r.DB, query, banner)
	if err != nil {
		return 0, err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)

	return dbtk.insertErrorHandler(ctx, result, err)
}

func (r *bannerRepo) Update(ctx context.Context, banner *models.Banner) error {
	query := `update banners set 
	slogan=:slogan, link_type=:link_type, 
	link_url=:link_url, image_url=:image_url, 
	enable=:enable, sort=:sort
	where id=:id and deleted_at is null;`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query, args, err := dbtk.debugSQL(ctx, r.DB, query, banner)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)

	return dbtk.updateErrorHandler(ctx, result, err)
}

func (r *bannerRepo) Get(ctx context.Context, id uint64) (banners *models.Banner, err error) {
	sql := `select * from banners where id = ? and deleted_at is null;`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	banners = new(models.Banner)

	err = r.DB.GetContext(ctx, banners, query, id)
	return banners, err
}

func (r *bannerRepo) List(ctx context.Context) (list []*models.Banner, err error) {
	sql := `select * from banners where deleted_at is null;`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query)

	err = r.DB.SelectContext(ctx, &list, query)
	return list, err
}

func (r *bannerRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update banners set deleted_at = now() where id = ?;`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)

	slog.DebugContext(ctx, query, slog.Int64("id", int64(id)))

	result, err := r.DB.ExecContext(ctx, query, id)
	return dbtk.updateErrorHandler(ctx, result, err)
}
