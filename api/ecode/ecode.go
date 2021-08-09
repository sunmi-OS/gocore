package ecode

import (
	"github.com/spf13/cast"
)

var errorMap = map[string]int{}

// New new error code and msg
func New(code int, err error) error {
	errorMap[err.Error()] = code
	return err
}

func Transform(err error) int {
	v, ok := errorMap[err.Error()]
	if !ok {
		return -1
	}
	return cast.ToInt(v)
}
