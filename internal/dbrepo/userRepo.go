package dbrepo

import (
	"context"
	"fmt"

	"github.com/lightsaid/ebook/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) (uint64, error)
	Delete(ctx context.Context, userID uint64) error
	Update(ctx context.Context, user *models.User) error
	Get(ctx context.Context, userID uint64) (*models.User, error)
	List(ctx context.Context, filter Filters) (*PageQueryVo, error)
}

var _ UserRepo = (*userRepo)(nil)

type userRepo struct {
	DB Queryable
}

func NewUserRepo(db Queryable) *userRepo {
	repo := &userRepo{
		DB: db,
	}

	return repo
}

func (r *userRepo) Create(ctx context.Context, user *models.User) (uint64, error) {
	query := `insert into users(
		email,
		password,
		nickname,
		avatar
	) values(
	 	:email,
		:password,
		:nickname,
		:avatar
	)`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query, arg, err := dbtk.debugSQL(ctx, r.DB, query, user)
	if err != nil {
		return 0, err
	}

	result, err := r.DB.ExecContext(ctx, query, arg...)
	return insertErrorHandler(result, err)
}

func (r *userRepo) Delete(ctx context.Context, userID uint64) error {
	query := `update users set deleted_at=now() where id=?`
	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query = r.DB.Rebind(query)

	return updateErrorHandler(r.DB.ExecContext(ctx, query, userID))
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `update users set 
		password = :password,
		nickname = :nickname,
		avatar = :avatar,
		role = :role,
		login_at = :login_at,
		login_ip = :login_ip
		where id = :id and deleted_at is null;
	`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query, arg, err := dbtk.debugSQL(ctx, r.DB, query, user)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query, arg...)
	return updateErrorHandler(result, err)
}

func (r *userRepo) Get(ctx context.Context, userID uint64) (*models.User, error) {
	query := r.DB.Rebind(`select * from users where id = ?`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	user := new(models.User)
	err := r.DB.GetContext(ctx, user, query, userID)
	return user, err
}

func (r *userRepo) List(ctx context.Context, filter Filters) (*PageQueryVo, error) {
	query := fmt.Sprintf(`
	select 
		id,
		email,
		nickname,
		avatar,
		role,
		login_at,
		login_ip,
		created_at,
		updated_at
	from users where deleted_at is null
	order by %s limit ? offset ?
	`, filter.sortColumnWithDefault(r))

	query = r.DB.Rebind(query)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	list := make([]*models.User, 0)
	stmt, err := r.DB.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.SelectContext(ctx, &list, query, filter.limit(), filter.offset())
	if err != nil {
		return nil, err
	}

	totalQuery := `select count(*) as total from users where deleted_at is null`
	stmt, err = r.DB.PreparexContext(ctx, totalQuery)
	if err != nil {
		return nil, err
	}

	var total = 0
	err = stmt.Get(totalQuery, &total)
	if err != nil {
		return nil, err
	}

	metadata := calculateMetadata(total, filter.PageNum, filter.PageSize)

	vo := makePageQueryVo(metadata, list)

	return vo, nil
}

func (r *userRepo) defaultSortSafelist() []string {
	return []string{
		"id", "email", "nickname", "role", "login_at", "created_at", "updated_at",
		"-id", "-email", "-nickname", "-role", "-login_at", "-created_at", "-updated_at",
	}
}
