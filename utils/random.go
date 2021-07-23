package utils

import (
	"math/rand"
	"time"

	"github.com/spf13/cast"
)

// RandomI 随机数
func RandomI(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// RandomI64 随机数
func RandomI64(min int, max int) int64 {
	return cast.ToInt64(RandomI(min, max))
}

// Random0Z 随机字符串 0~Z
func Random0Z(l int) string {
	return Random(l, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// Random09 随机字符串 0-9
func Random09(l int) string {
	return Random(l, "0123456789")
}

// Random 随机字符串指定范围
func Random(l int, randomRange string) string {
	bytes := []byte(randomRange)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
