package dbcache

import "github.com/redis/go-redis/v9"

var (
	rdb *redis.Client
)

type Repository struct {
	UserCache UserCache
	BookCache BookCache
}

func NewRepository(client *redis.Client) Repository {
	// 保存到包下使用
	rdb = client

	return Repository{
		UserCache: NewUserCache(),
		BookCache: NewBookCache(),
	}
}
