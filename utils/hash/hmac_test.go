package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHmacSha1(t *testing.T) {
	s, _ := HmacSha1("gocore", "123")
	assert.Equal(t, s, "e86ea7bb3aac9d45c2942b653c302d94f3349656")
}

func TestHmacSha224(t *testing.T) {
	s, _ := HmacSha224("gocore", "123")
	assert.Equal(t, s, "181981ed939dd37895d3382f050e51859203fffba35aa7ff4ed97bd4")
}

func TestHmacSha256(t *testing.T) {
	s, _ := HmacSha256("gocore", "123")
	assert.Equal(t, s, "ba2f4a58fd234d3a04ad31107b0a44e4345de83a93eea63a0a28ec521c88af81")
}

func TestHmacSha384(t *testing.T) {
	s, _ := HmacSha384("gocore", "123")
	assert.Equal(t, s, "d783d3dd8aaddc808a7cbcc697cdd238cf82d203923053a1bc296b73219dadef8f091c20276a0b3bf8ba5b3c1f16b785")
}

func TestHmacSha512(t *testing.T) {
	s, _ := HmacSha512("gocore", "123")
	assert.Equal(t, s, "f5c9a05b860c83d6084ad6f3ec227e8f3d2b5bd072f94f0904e0d8de514d826816feca8999ca230074cd3c6cd49e0bca472914cf3efd1ae1a39ccef4e8917bbd")
}

func TestHmacMD5(t *testing.T) {
	s, _ := HmacMD5("gocore", "123")
	assert.Equal(t, s, "a1f6b37b45a1edfd420fa9be02f05b28")
}
