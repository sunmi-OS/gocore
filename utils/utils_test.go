package utils

import (
	"log"
	"testing"
)

func TestRandInt(t *testing.T) {
	randInt := RandInt(100, 120)
	log.Println(randInt)
}
