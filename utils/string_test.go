package utils

import "testing"

func TestJoinInts(t *testing.T) {
	tests := []struct {
		name string
		is   []int64
		want string
	}{
		{"test1", []int64{1, 2, 3, 4, 5}, "1,2,3,4,5"},
		{"test2", []int64{1}, "1"},
		{"test3", []int64{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinInts(tt.is); got != tt.want {
				t.Errorf("JoinInts() = %v, want %v", got, tt.want)
			}
		})
	}
}
