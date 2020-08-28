package alog

import "testing"

func TestInit(t *testing.T) {
	c := &LoggerConfig{
		ConfigName: "TestLog",
		LogStore:   "TestLogStore",
	}
	New(c)
}
