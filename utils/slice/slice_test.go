package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	var validTests = []struct {
		s1 []int
		s2 []int
		s3 []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{4, 5, 6}, []int{1, 2, 3}},
		{[]int{1, 3, 5}, []int{2, 4, 7}, []int{1, 3, 5}},
	}

	for _, tt := range validTests {
		assert.Equal(t, tt.s3, Diff(tt.s1, tt.s2))
	}
}

func TestIntersect(t *testing.T) {
	var validTests = []struct {
		s1 []int
		s2 []int
		s3 []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{4, 5, 6}, []int{4, 5}},
		{[]int{1, 3, 5}, []int{1, 2, 4, 7}, []int{1}},
	}

	for _, tt := range validTests {
		assert.Equal(t, tt.s3, Intersect(tt.s1, tt.s2))
	}
}

var compactTests = []struct {
	name string
	s    []int
	want []int
}{
	{
		"nil",
		nil,
		nil,
	},
	{
		"one",
		[]int{1},
		[]int{1},
	},
	{
		"sorted",
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		"1 item",
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		"unsorted",
		[]int{1, 2, 1},
		[]int{1, 2, 1},
	},
	{
		"many",
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestCompact(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := Compact(copy); !Equal(got, test.want) {
			t.Errorf("Compact(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}
