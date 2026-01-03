package dbcache

type UserCache interface{}

type userCache struct {
}

var _ UserCache = (*userCache)(nil)

func NewUserCache() *userCache {
	return &userCache{}
}
