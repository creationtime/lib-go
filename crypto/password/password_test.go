package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	pwd    = "testPassword$x1w432b7^"
	hashed = "argon2id$1$65536$12$32$//mhs71hH9GtEN721b2kN2+1RmiUZBZvgabYXl0xaxk=$mKgtftWvjvzHQMqOCv6NrjgvjfudP4/lc/dAyC66giw="
)

func TestGenerateSaltedHash(t *testing.T) {
	v, err := GenerateSaltedHash(pwd)
	assert.NoError(t, err)
	assert.NotEmpty(t, v)
}

func BenchmarkGenerateSaltedHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, err := GenerateSaltedHash(pwd)
		assert.NoError(b, err)
		assert.NotEmpty(b, v)
	}
}

func TestCompareHashWithPassword(t *testing.T) {
	ok, err := CompareHashWithPassword(hashed, pwd)
	assert.NoError(t, err)
	assert.EqualValues(t, true, ok)
}

func BenchmarkCompareHashWithPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ok, err := CompareHashWithPassword(hashed, pwd)
		assert.NoError(b, err)
		assert.EqualValues(b, true, ok)
	}
}
