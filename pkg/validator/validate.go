package validator

import "regexp"

// 定义业务字段验证

var (
	isbn13 = regexp.MustCompile("^(?:[0-9]{13})$")
	isbn10 = regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$")
)

// IsISBN 检查是否是 isbn
func IsISBN(isbn string) bool {
	if isbn10.MatchString(isbn) || isbn13.MatchString(isbn) {
		return true
	}
	return false
}
