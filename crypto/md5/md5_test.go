package md5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	s = "aaa"
	p = "47bce5c74f589f4867dbd57e9ca9f808"
)

func TestHash(t *testing.T) {
	v := Hash(s)
	assert.EqualValues(t, p, v)
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := Hash(s)
		assert.EqualValues(b, p, v)
	}
}
