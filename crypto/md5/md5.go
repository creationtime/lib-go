package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
