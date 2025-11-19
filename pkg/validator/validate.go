package validator

import (
	"regexp"
)

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

// Verifiyer定义一个Verifiy校验接口方法，
// 配合DoVerifiy使用可以省去validator.New(),Valid()步骤
type Verifiyer interface {
	Verifiy(*Validator)
}

// DoVerifiy 执行实现Verifiyer接口对象，并返回*Validator
func DoVerifiy(obj Verifiyer) (*Validator, bool) {
	v := New()
	obj.Verifiy(v)
	return v, v.Valid()
}
