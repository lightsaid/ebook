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
	GetByUqField(ctx context.Context, uq UserUq) (*models.User, error)
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
	return dbtk.insertErrorHandler(ctx, result, err)
}

func (r *userRepo) Delete(ctx context.Context, userID uint64) error {
	query := r.DB.Rebind(`update users set deleted_at=now() where id=?`)
	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, userID)
	return dbtk.updateErrorHandler(ctx, result, err)
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
	return dbtk.updateErrorHandler(ctx, result, err)
}

func (r *userRepo) Get(ctx context.Context, userID uint64) (*models.User, error) {
	query := r.DB.Rebind(`select * from users where id = ? and deleted_at is null;`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	user := new(models.User)
	err := r.DB.GetContext(ctx, user, query, userID)
	return user, err
}

type UserUq struct {
	ID    uint64
	Email string
}

// GetByUqField 根据唯一字段查询单个用户信息 uq 提供id或email即可，
// 如果两者都不提供，则返回 ErrNotFound
func (r *userRepo) GetByUqField(ctx context.Context, uq UserUq) (*models.User, error) {
	if uq.ID > 0 {
		return r.Get(ctx, uq.ID)
	}

	if uq.Email == "" {
		return nil, ErrNotFound
	}
	query := r.DB.Rebind(`select * from users where email = ? and deleted_at is null;`)

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	user := new(models.User)
	err := r.DB.GetContext(ctx, user, query, uq.Email)
	return user, err
}

func (r *userRepo) List(ctx context.Context, filter Filters) (*PageQueryVo, error) {
	// 检查分页设置
	filter.check()
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

	err := r.DB.SelectContext(ctx, &list, query, filter.limit(), filter.offset())
	if err != nil {
		return nil, err
	}

	totalQuery := `select count(*) as total from users where deleted_at is null`

	var total = 0
	err = r.DB.GetContext(ctx, &total, totalQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println("========", filter)
	metadata := dbtk.calculateMetadata(total, filter.PageNum, filter.PageSize)

	vo := dbtk.makePageQueryVo(metadata, list)

	return vo, nil
}

func (r *userRepo) defaultSortSafelist() []string {
	return []string{
		"id", "email", "nickname", "role", "login_at", "created_at", "updated_at",
		"-id", "-email", "-nickname", "-role", "-login_at", "-created_at", "-updated_at",
	}
}
