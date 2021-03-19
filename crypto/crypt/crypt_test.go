package crypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	key     = "dONkmXxhl6kLRiLt"
	iv      = "dONkmXxhl6kLRiLt"
	mode    = ModeCBC
	padding = "pkcs5"
	coding  = "base64"

	plainText  = "test"
	cipherText = "2Um4/PuuXqABtDAQEbFrW2n60cu7KE0IGACyyFKdEnukoER7Zi0PO0xJzsscarGdSadm1QCiejKRJqsoNZ6YYtdi9Z2Jkl5Ey+fQZZcfV9hIf4MsWdOHb5dfPU/c6q13"
)

func TestNew(t *testing.T) {
	c, err := New("aes", key, iv, mode, padding, coding)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestCrypto_Encode(t *testing.T) {
	c, err := New("aes", key, iv, mode, padding, coding)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	content, err := c.Encode(plainText)
	assert.NoError(t, err)
	assert.EqualValues(t, cipherText, content)
}

func BenchmarkCrypto_Encode(b *testing.B) {
	c, err := New("aes", key, iv, mode, padding, coding)
	assert.NoError(b, err)
	assert.NotNil(b, c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Encode(plainText)
		assert.NoError(b, err)
	}
}

func TestCrypto_Decode(t *testing.T) {
	c, err := New("aes", key, iv, mode, padding, coding)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	content, err := c.Decode(cipherText)
	assert.NoError(t, err)
	assert.EqualValues(t, plainText, content)
}

func BenchmarkCrypto_Decode(b *testing.B) {
	c, err := New("aes", key, iv, mode, padding, coding)
	assert.NoError(b, err)
	assert.NotNil(b, c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Decode(cipherText)
		assert.NoError(b, err)
	}
}
