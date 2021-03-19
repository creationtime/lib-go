package pad

import (
	"bytes"
	"errors"
)

type Padding interface {
	Pad(cipherText []byte) []byte
	UnPad(origData []byte) ([]byte, error)
}

type Zero struct {
	BlockSize int
}

func NewZero(bs int) Padding {
	return &Zero{BlockSize: bs}
}

func (z *Zero) Pad(cipherText []byte) []byte {
	padding := z.BlockSize - len(cipherText)%z.BlockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padText...)
}

func (z *Zero) UnPad(origData []byte) ([]byte, error) {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		}), nil
}

type PKCS5 struct {
	BlockSize int
}

func NewPKCS5(bs int) Padding {
	return &PKCS5{BlockSize: bs}
}

func (p *PKCS5) Pad(cipherText []byte) []byte {
	padding := p.BlockSize - len(cipherText)%p.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (p *PKCS5) UnPad(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if length < unPadding {
		return nil, errors.New("unpadding error")
	}
	return origData[:(length - unPadding)], nil
}

type PKCS7 struct {
	BlockSize int
}

func NewPKCS7(bs int) Padding {
	return &PKCS7{BlockSize: bs}
}

func (p *PKCS7) Pad(cipherText []byte) []byte {
	padding := p.BlockSize - len(cipherText)%p.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (p *PKCS7) UnPad(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if length < unPadding {
		return nil, errors.New("unpadding error")
	}
	return origData[:(length - unPadding)], nil
}
