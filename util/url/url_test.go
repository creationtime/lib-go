package url

import (
	"gotest.tools/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	baseUrl := "http://www.baidu.com"
	t.Run("https", func(t *testing.T) {
		s := Join("https://blog.csdn.net/qq_27825451/article/details/89309175", baseUrl)
		assert.Equal(t, s, "https://blog.csdn.net/qq_27825451/article/details/89309175")
	})
	t.Run("http", func(t *testing.T) {
		s := Join("http://blog.csdn.net/qq_27825451/article/details/89309175", baseUrl)
		assert.Equal(t, s, "http://blog.csdn.net/qq_27825451/article/details/89309175")
	})
	t.Run("//", func(t *testing.T) {
		s := Join("//blog.csdn.net/qq_27825451/article/details/89309175", baseUrl)
		assert.Equal(t, s, "//blog.csdn.net/qq_27825451/article/details/89309175")
	})
	t.Run("/", func(t *testing.T) {
		s := Join("/blog.csdn.net/qq_27825451/article/details/89309175", baseUrl)
		assert.Equal(t, s, "http://www.baidu.com/blog.csdn.net/qq_27825451/article/details/89309175")
	})
}

func TestTrim(t *testing.T) {
	t.Run("https_query", func(t *testing.T) {
		s := Trim("https://blog.csdn.net/qq_27825451/article/details/89309175?a=1")
		assert.Equal(t, s, "/qq_27825451/article/details/89309175?a=1")
	})
	t.Run("https", func(t *testing.T) {
		s := Trim("https://blog.csdn.net/qq_27825451/article/details/89309175")
		assert.Equal(t, s, "/qq_27825451/article/details/89309175")
	})
	t.Run("http", func(t *testing.T) {
		s := Trim("http://blog.csdn.net/qq_27825451/article/details/89309175")
		assert.Equal(t, s, "/qq_27825451/article/details/89309175")
	})
	t.Run("//", func(t *testing.T) {
		s := Trim("//blog.csdn.net/qq_27825451/article/details/89309175")
		assert.Equal(t, s, "/qq_27825451/article/details/89309175")
	})
	t.Run("/", func(t *testing.T) {
		s := Trim("/blog.csdn.net/qq_27825451/article/details/89309175")
		assert.Equal(t, s, "/blog.csdn.net/qq_27825451/article/details/89309175")
	})
}
