package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, Uniq([]int{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []int64{1, 2, 3, 4}, Uniq([]int64{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []float64{1.01, 2.02, 3.03, 4.04}, Uniq([]float64{1.01, 2.02, 1.01, 3.03, 4.04, 3.03}))
	assert.Equal(t, []string{"h", "e", "l", "o"}, Uniq([]string{"h", "e", "l", "l", "o"}))
}

func TestWithout(t *testing.T) {
	result1 := Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
	result2 := Without([]int{0, 7}, 0, 1, 2, 3, 4, 5)
	result3 := Without([]int{}, 0, 1, 2, 3, 4, 5)
	result4 := Without([]int{0, 1, 2}, 0, 1, 2)
	result5 := Without([]int{})
	assert.Equal(t, []int{10}, result1)
	assert.Equal(t, []int{7}, result2)
	assert.Equal(t, []int{}, result3)
	assert.Equal(t, []int{}, result4)
	assert.Equal(t, []int{}, result5)
}

func TestUnion(t *testing.T) {
	result1 := Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
	result2 := Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7})
	result3 := Union([]int{0, 1, 2, 3, 4, 5}, []int{})
	result4 := Union([]int{0, 1, 2}, []int{0, 1, 2})
	result5 := Union([]int{}, []int{})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 10}, result1)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, result2)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, result3)
	assert.Equal(t, []int{0, 1, 2}, result4)
	assert.Equal(t, []int{}, result5)

	result11 := Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10}, []int{0, 1, 11})
	result12 := Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7}, []int{8, 9})
	result13 := Union([]int{0, 1, 2, 3, 4, 5}, []int{}, []int{})
	result14 := Union([]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2})
	result15 := Union([]int{}, []int{}, []int{})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 10, 11}, result11)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, result12)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, result13)
	assert.Equal(t, []int{0, 1, 2}, result14)
	assert.Equal(t, []int{}, result15)
}

func TestRand(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret1 := Rand(a1, 2)
	assert.Equal(t, 2, len(ret1))
	assert.NotEqual(t, a1[:2], ret1)

	a2 := []float64{1.01, 2.02, 3.03, 4.04, 5.05, 6.06, 7.07, 8.08, 9.09, 10.10}
	ret2 := Rand(a2, 2)
	assert.Equal(t, 2, len(ret2))
	assert.NotEqual(t, a2[:2], ret2)

	a3 := []string{"h", "e", "l", "l", "o", "w", "o", "r", "l", "d"}
	ret3 := Rand(a3, 2)
	assert.Equal(t, 2, len(ret3))
	assert.NotEqual(t, a3[:2], ret3)

	type User struct {
		ID   int64
		Name string
	}

	a4 := []User{
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   4,
			Name: "l",
		},
		{
			ID:   5,
			Name: "o",
		},
		{
			ID:   6,
			Name: "w",
		},
		{
			ID:   7,
			Name: "o",
		},
		{
			ID:   8,
			Name: "r",
		},
		{
			ID:   9,
			Name: "l",
		},
		{
			ID:   10,
			Name: "d",
		},
	}

	ret4 := Rand(a4, 2)
	assert.Equal(t, 2, len(ret4))
	assert.NotEqual(t, a4[:2], ret4)

	ret5 := Rand(a4, -1)
	assert.Equal(t, len(a4), len(ret5))
	assert.NotEqual(t, a4, ret5)
}

func TestPinTop(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 3)
	assert.Equal(t, []int{4, 1, 2, 3, 5}, a1)

	a2 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 0)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a2)

	a3 := []int{1, 2, 3, 4, 5}
	PinTop(a1, -1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a3)

	a4 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a4)
}

func TestPinTopF(t *testing.T) {
	type Demo struct {
		ID   int
		Name string
	}
	arr := []Demo{
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   4,
			Name: "o",
		},
		{
			ID:   5,
			Name: "w",
		},
	}
	PinTopF(arr, func(v Demo) bool {
		return v.Name == "o"
	})
	assert.Equal(t, []Demo{
		{
			ID:   4,
			Name: "o",
		},
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   5,
			Name: "w",
		},
	}, arr)
}

func TestChunk(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ret1 := Chunk(a, 2)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}, ret1)

	ret2 := Chunk(a, 3)
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, ret2)

	ret3 := Chunk(a, 4)
	assert.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}, ret3)
}
