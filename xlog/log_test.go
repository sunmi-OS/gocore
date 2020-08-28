package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {

	Error(map[string]interface{}{
		"name": "jerry",
		"age":  14,
	})

	Info(struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	})

	Warn("Warn")
	Debug("Debug")
	Error("Error")
}
