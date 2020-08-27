package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {

	Error(map[string]interface{}{
		"name": "jerry",
		"age":  14,
	})
	Info("info")
	Warning("warning")
}
