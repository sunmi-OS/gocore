package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha1(t *testing.T) {
	s, _ := Sha1("gocore")
	assert.Equal(t, s, "5e4f3bbfb45318ddb74e86f772fbe3626132ab20")
}

func TestSha224(t *testing.T) {
	s, _ := Sha224("gocore")
	assert.Equal(t, s, "3042a8685b4d129711ab55ee743f87e47940e3385d9ec8ece50cd935")
}

func TestSha256(t *testing.T) {
	s, _ := Sha256("gocore")
	assert.Equal(t, s, "ad1456cc07aa133df96689cba47eb7ce3b8a41d5f0997034400e9d049a40f4ab")
}

func TestSha384(t *testing.T) {
	s, _ := Sha384("gocore")
	assert.Equal(t, s, "d73ae09b1e5704f8cd64962636266a20ccb4c78927a7bd42e1e1589696d4299a5c87f1cfe7ebe17ad62688c03cda96aa")
}

func TestSha512(t *testing.T) {
	s, _ := Sha512("gocore")
	assert.Equal(t, s, "89a9a722c4f913493a925808f11c31236444315cd2438396764e908acc394d593445915d222e36094b68e7dc2266b75f4a32ce8a7e20dd7edc5caa19fc3fccc6")
}
