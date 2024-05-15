// Package iterator provides convenient functionality for batch execution.
// Author: Jerry
package iterator

type Step struct {
	Head int
	Tail int
}

// Steps calculates the steps.
// example:
//
//	ids := []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
//	for _, step := range iterator.Steps(len(ids), 10) {
//		cur := ids[step.Head:step.Tail]
//		// todo: do something
//	}
func Steps(total, step int) (steps []Step) {
	steps = make([]Step, 0)
	for i := 0; i < total; i++ {
		if i%step == 0 {
			head := i
			tail := head + step
			if tail > total {
				tail = total
			}
			steps = append(steps, Step{Head: head, Tail: tail})
		}
	}
	return steps
}
