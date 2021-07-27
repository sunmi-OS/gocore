package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	assert.Equal(t, MD5("gocore"), "137c5840f84e21a0e55aa884f971a166")
}

func TestMD5File(t *testing.T) {
	s, _ := MD5File("./test")
	assert.Equal(t, s, "137c5840f84e21a0e55aa884f971a166")
}
