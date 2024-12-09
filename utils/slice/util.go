package slice

import (
	"math"
	"math/rand/v2"
)

// Uniq 集合去重
func Uniq[T ~[]E, E comparable](list T) T {
	if len(list) == 0 {
		return list
	}

	ret := make(T, 0, len(list))
	m := make(map[E]struct{}, len(list))
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// Without 返回不包括所有给定值的切片
func Without[T ~[]E, E comparable](list T, exclude ...E) T {
	if len(list) == 0 {
		return list
	}

	m := make(map[E]struct{}, len(exclude))
	for _, v := range exclude {
		m[v] = struct{}{}
	}

	ret := make(T, 0, len(list))
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// Union 返回两个集合的并集
func Union[T ~[]E, E comparable](lists ...T) T {
	ret := make(T, 0)
	m := make(map[E]struct{})
	for _, list := range lists {
		for _, v := range list {
			if _, ok := m[v]; !ok {
				ret = append(ret, v)
				m[v] = struct{}{}
			}
		}
	}
	return ret
}

// Rand 返回一个指定随机挑选个数的切片
// 若 n == -1 or n >= len(list)，则返回打乱的切片
func Rand[T ~[]E, E any](list T, n int) T {
	if n == 0 || n < -1 {
		return nil
	}

	count := len(list)
	ret := make(T, count)
	copy(ret, list)
	rand.Shuffle(count, func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	if n == -1 || n >= count {
		return ret
	}
	return ret[:n]
}

// PinTop 置顶集合中的一个元素
func PinTop[T ~[]E, E any](list T, index int) {
	if index <= 0 || index >= len(list) {
		return
	}
	for i := index; i > 0; i-- {
		list[i], list[i-1] = list[i-1], list[i]
	}
}

// PinTopF 置顶集合中满足条件的一个元素
func PinTopF[T ~[]E, E any](list T, fn func(v E) bool) {
	index := 0
	for i, v := range list {
		if fn(v) {
			index = i
			break
		}
	}
	for i := index; i > 0; i-- {
		list[i], list[i-1] = list[i-1], list[i]
	}
}

// Chunk 集合分片
func Chunk[T ~[]E, E any](list T, size int) []T {
	if size <= 0 {
		return []T{}
	}
	length := len(list)
	count := int(math.Ceil(float64(length) / float64(size)))
	ret := make([]T, 0, count)
	for i := 0; i < count; i++ {
		start := i * size
		end := (i + 1) * size
		if end > length {
			end = length
		}
		ret = append(ret, list[start:end])
	}
	return ret
}
