package dbrepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lightsaid/ebook/internal/models"
)

type AuthorRepo interface {
	Create(ctx context.Context, authorName string) (uint64, error)
	Update(ctx context.Context, id uint64, authorName string) error
	Get(ctx context.Context, id uint64) (*models.Author, error)
	List(ctx context.Context, f Filters) (*PageQueryVo, error)
	Delete(ctx context.Context, id uint64) error
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

func (r *authorRepo) Create(ctx context.Context, authorName string) (uint64, error) {
	sql := `insert author(author_name) values(:author_name);`

	ctx, cancal := timeoutCtx(ctx)
	defer cancal()

	query, args, err := debugSQL(ctx, r.DB, sql, envelop{"author_name": authorName})
	if err != nil {
		return 0, err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	return insertErrorHandler(result, err)
}

func (r *authorRepo) Update(ctx context.Context, id uint64, authorName string) error {
	sql := `update author set author_name=:author_name where id=:id and deleted_at is null;`

	ctx, cancal := timeoutCtx(ctx)
	defer cancal()

	query, args, err := debugSQL(
		ctx, r.DB, sql,
		envelop{"author_name": authorName, "id": id},
	)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
	return updateErrorHandler(result, err)
}

func (r *authorRepo) Get(ctx context.Context, id uint64) (author *models.Author, err error) {
	sql := `
		select 
			id, author_name, created_at, updated_at 
		from 
			author 
		where 
			id=? and deleted_at is null;`
	author = new(models.Author)

	ctx, cancal := timeoutCtx(ctx)
	defer cancal()

	slog.InfoContext(ctx, sql, slog.Int64("id", int64(id)))

	err = r.DB.GetContext(ctx, author, sql, id)
	return author, err
}

// List 分页获取
func (r *authorRepo) List(ctx context.Context, f Filters) (*PageQueryVo, error) {
	// 如果没有使用默认的
	if len(f.SortSafelist) == 0 {
		f.SortSafelist = r.defaultSortSafelist()
	}

	sortFields := f.sortColumn()

	// 提供默认
	if sortFields == "" {
		sortFields = "id ASC"
	}

	sql := fmt.Sprintf(
		`select 
			*, 
			(select count(*) from author where deleted_at is null) as total
		 from author where deleted_at is null
		 order by %s
		 limit ? offset ?`, sortFields,
	)

	ctx, cancal := timeoutCtx(ctx)
	defer cancal()

	query := r.DB.Rebind(sql)

	slog.InfoContext(
		ctx,
		spaceRex.ReplaceAllString(query, " "),
		slog.String("sort", sortFields),
		slog.Int("limit", f.limit()),
		slog.Int("offset", f.offset()),
	)

	type SQLAutor struct {
		models.Author
		Total int `db:"total"`
	}

	arr := make([]*SQLAutor, 0)

	err := r.DB.SelectContext(ctx, &arr, query, f.limit(), f.offset())
	if err != nil {
		return nil, err
	}

	list := make([]*models.Author, 0, len(arr))

	for _, x := range arr {
		list = append(list, &x.Author)
	}

	var metadata Metadata
	if len(arr) > 0 {
		metadata = calculateMetadata(arr[0].Total, f.PageNum, f.PageSize)
	}

	vo := makePageQueryVo(metadata, list)

	return vo, err
}

// Delete 删除图书
func (r *authorRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update author set deleted_at = now() where id = ?;`

	ctx, cancal := timeoutCtx(ctx)
	defer cancal()

	slog.InfoContext(ctx, sql, slog.Int64("id", int64(id)))

	result, err := r.DB.ExecContext(ctx, sql, id)
	return updateErrorHandler(result, err)
}

// defaultSortSafelist 导出默认的安全排序字段
func (r *authorRepo) defaultSortSafelist() []string {
	return []string{
		"id", "author_name", "created_at", "updated_at",
		"-id", "-author_name", "-created_at", "-updated_at",
	}
}
