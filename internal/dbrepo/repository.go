package dbrepo

type Repository struct {
	BookRepo         BookRepo
	AuthorRepo       AuthorRepo
	CategoryRepo     CategoryRepo
	PublisherRepo    PublisherRepo
	BookCategoryRepo BookCategoryRepo
	BannerRepo       BannerRepo
	UserRepo         UserRepo
	OrderRepo        OrderRepo
	ShoppingCartRepo ShoppingCartRepo
}

// NewRepository创建一个Repository仓库，使用Queryable接口，同时兼容sql.DB和sql.Tx方法
func NewRepository(db Queryable) Repository {
	return Repository{
		BookRepo:         NewBookRepo(db),
		AuthorRepo:       NewAuthorRepo(db),
		CategoryRepo:     NewCategoryRepo(db),
		PublisherRepo:    NewPublisherRepo(db),
		BookCategoryRepo: NewBookCategoryRepo(db),
		BannerRepo:       NewBannerRepo(db),
		UserRepo:         NewUserRepo(db),
		OrderRepo:        NewOrderRepo(db),
		ShoppingCartRepo: NewShoppingCartRepo(db),
	}
}
