package encoding

import (
	"encoding/base64"
	"encoding/hex"
)

type Encoding interface {
	Encode(src []byte) string
	Decode(src string) ([]byte, error)
}

type Hex struct {
}

func NewHex() Encoding {
	return &Hex{}
}

func (h *Hex) Encode(src []byte) string {
	return hex.EncodeToString(src)
}

func (h *Hex) Decode(src string) ([]byte, error) {
	return hex.DecodeString(src)
}

type Base64 struct {
}

func NewBase64() Encoding {
	return &Base64{}
}

func (b *Base64) Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func (b *Base64) Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

type Base64URL struct {
}

func NewBase64URL() Encoding {
	return &Base64URL{}
}

func (b *Base64URL) Encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func (b *Base64URL) Decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
