package utils

import (
	"log"
	"testing"
)

func TestSteps(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for index, step := range Steps(len(arr), 6) {
		ids := arr[step.Head:step.Tail]
		log.Printf("step[%d], slice: %d\n", index, ids)
	}
}
