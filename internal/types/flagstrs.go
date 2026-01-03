package types

import "fmt"

// 例如：go run main.go -env "develop.env" -env "api.develop.env"

// ArrayString 自定义一个类型，实现flag解析自定义字符串数组
type ArrayString []string

// String 实现 String() 方法
func (s *ArrayString) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set 实现 Set() 方法
func (s *ArrayString) Set(value string) error {
	*s = append(*s, value)
	return nil
}
