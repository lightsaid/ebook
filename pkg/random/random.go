package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const characters = "1234567890qwertyuiopasdfghjklzxcvbnmWERTYUIOPASDFGHJKLZXCVBNM"

func init() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成随机整数，在min和max之间
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomString 生成随机字符串
func RandomString(n int) string {
	// 声明一个字符串构造器
	var sb strings.Builder
	k := len(characters)
	for i := 0; i < n; i++ {
		c := characters[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// UUID 生成一个简单uuid
func UUID() string {
	return fmt.Sprintf("%d%s", time.Now().UnixMicro(), RandomString(6))
}
