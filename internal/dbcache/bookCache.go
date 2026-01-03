package dbcache

type BookCache interface{}

type bookCache struct {
}

var _ BookCache = (*bookCache)(nil)

func NewBookCache() *bookCache {
	return &bookCache{}
}
