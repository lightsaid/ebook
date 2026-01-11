package dbcache

import "fmt"

const (
	baseKey = "ebook:admin"
)

func userIDKey(id uint64) string {
	return fmt.Sprintf("%s:user:%d", baseKey, id)
}
