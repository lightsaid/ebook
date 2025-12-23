package dbrepo

type OrderRepo interface{}

var _ OrderRepo = (*orderRepo)(nil)

type orderRepo struct {
	DB Queryable
}

func NewOrderRepo(db Queryable) *orderRepo {
	repo := &orderRepo{
		DB: db,
	}

	return repo
}
