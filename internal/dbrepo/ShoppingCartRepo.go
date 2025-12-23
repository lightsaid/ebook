package dbrepo

import (
	"context"
	"fmt"

	"github.com/lightsaid/ebook/internal/models"
)

type ShoppingCartRepo interface {
	Create(ctx context.Context, data *models.ShoppingCart) (uint64, error)
	Delete(ctx context.Context, shoppingId uint64) error
	List(ctx context.Context, userId uint64, f Filters) (*PageQueryVo, error)
}

var _ ShoppingCartRepo = (*shoppingCartRepo)(nil)

type shoppingCartRepo struct {
	DB Queryable
}

func NewShoppingCartRepo(db Queryable) *shoppingCartRepo {
	repo := &shoppingCartRepo{
		DB: db,
	}

	return repo
}

func (r *shoppingCartRepo) Create(ctx context.Context, data *models.ShoppingCart) (uint64, error) {
	query := "insert into shopping_carts(user_id,book_id,quantity) values(:user_id,:book_id,:quantity)"
	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	query, args, err := dbtk.debugSQL(ctx, r.DB, query, data)
	if err != nil {
		return 0, err
	}

	stmt, err := r.DB.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	return insertErrorHandler(stmt.ExecContext(ctx, args...))
}

func (r *shoppingCartRepo) Delete(ctx context.Context, shoppingId uint64) error {
	query := `delete from shopping_carts where id=?`
	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	return updateErrorHandler(r.DB.ExecContext(ctx, query, shoppingId))
}

func (r *shoppingCartRepo) List(ctx context.Context, userId uint64, f Filters) (*PageQueryVo, error) {
	totalQuery := `select count(*) as total from shopping_carts sc 
	left join books b on sc.book_id = b.id where sc.user_id = ?`

	ctx, cancel := dbtk.withTimeout(ctx)
	defer cancel()

	var total int
	err := r.DB.GetContext(ctx, &total, totalQuery, userId)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
	select sc.*, b.* from shopping_carts sc 
	left join books b on sc.book_id = b.id where sc.user_id = ? 
	order by %s limit ? offset ?`, "sc.created_at DESC, sc.id ASC")

	stmt, err := r.DB.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var list = make([]*models.SQLShoppingCart, 0, f.PageSize)
	stmt.SelectContext(ctx, &list, userId, f.limit(), f.offset())

	metadata := calculateMetadata(total, f.PageNum, f.PageSize)
	vo := makePageQueryVo(metadata, list)

	return vo, nil
}

// TODO:
func (r *shoppingCartRepo) ListByBookName(ctx context.Context, userId uint64, f Filters) (*PageQueryVo, error) {
	return nil, nil
}
