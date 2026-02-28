package dbcache

import "fmt"

const (
	baseAdminKey  = "ebook:admin"
	basePortalKey = "ebook:portal"
)

func userIDKey(id uint64) string {
	return fmt.Sprintf("%s:user:%d", baseAdminKey, id)
}
