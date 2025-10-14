package random

import (
	"math/rand"
	"strings"
	"time"
)

var srcRand *rand.Rand

const characters = "1234567890qwertyuiopasdfghjklzxcvbnmWERTYUIOPASDFGHJKLZXCVBNM"

func init() {
	// 设置随机种子
	srcRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt 生成随机整数，在min和max之间
func RandomInt(min, max int) int {
	return min + srcRand.Intn(max-min+1)
}

// RandomString 生成随机字符串
func RandomString(n int) string {
	// 声明一个字符串构造器
	var sb strings.Builder
	size := len(characters)

	for range n {
		s := characters[srcRand.Intn(size)]
		sb.WriteByte(s)
	}
	return sb.String()
}
