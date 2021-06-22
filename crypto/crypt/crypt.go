package crypt

import (
	"errors"
	strings2 "strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"github.com/lights-T/lib-go/crypto/crypt/ecb"
	"github.com/lights-T/lib-go/crypto/encoding"
	"github.com/lights-T/lib-go/crypto/pad"
	"github.com/lights-T/lib-go/util/strings"
)

type Mode uint8

const (
	ModeCBC Mode = iota
	ModeECB      //
	ModeCFB
	ModeOFB
)

type Crypt interface {
	Encode(string) (string, error)
	Decode(string) (string, error)
}

type Crypto struct {
	key      []byte
	iv       []byte
	mode     Mode
	block    cipher.Block
	padding  pad.Padding
	encoding encoding.Encoding
}

func New(a string, key string, iv string, mode Mode, padding string, coding string) (Crypt, error) {
	crypt := &Crypto{
		key:  strings.StringToBytes(key),
		iv:   strings.StringToBytes(iv),
		mode: mode,
	}
	var (
		block cipher.Block
		err   error
	)

	switch strings2.ToLower(a) {
	case "aes":
		block, err = aes.NewCipher(crypt.key)
	case "des":
		block, err = des.NewCipher(crypt.key)
	case "3des":
		block, err = des.NewTripleDESCipher(crypt.key)
	default:
		err = errors.New("algorithm not supported")
	}
	if err != nil {
		return nil, err
	}

	crypt.block = block
	switch strings2.ToLower(padding) {
	case "zero":
		crypt.padding = pad.NewZero(block.BlockSize())
	case "pkcs5":
		crypt.padding = pad.NewPKCS5(block.BlockSize())
	case "pkcs7":
		crypt.padding = pad.NewPKCS7(block.BlockSize())
	default:
		return nil, errors.New("padding method not supported")
	}

	switch strings2.ToLower(coding) {
	case "hex":
		crypt.encoding = encoding.NewHex()
	case "base64":
		crypt.encoding = encoding.NewBase64()
	case "base64Url":
		crypt.encoding = encoding.NewBase64URL()
	default:
		return nil, errors.New("encoding method not supported")
	}
	return crypt, nil
}

func (c *Crypto) Encode(src string) (string, error) {
	plainText := c.padding.Pad(strings.StringToBytes(src))
	cipherText := make([]byte, len(plainText))

	switch c.mode {
	case ModeCBC:
		enc := cipher.NewCBCEncrypter(c.block, c.iv)
		enc.CryptBlocks(cipherText, plainText)
	case ModeECB:
		enc := ecb.NewECBEncrypter(c.block)
		enc.CryptBlocks(cipherText, plainText)
	case ModeCFB:
		enc := cipher.NewCFBEncrypter(c.block, c.iv)
		enc.XORKeyStream(cipherText, plainText)
	case ModeOFB:
		enc := cipher.NewOFB(c.block, c.iv)
		enc.XORKeyStream(cipherText, plainText)
	default:
		return "", errors.New("invalid cipher mode")
	}
	return c.encoding.Encode(cipherText), nil
}

func (c *Crypto) Decode(src string) (string, error) {
	cipherText, err := c.encoding.Decode(src)
	if err != nil {
		return "", err
	}
	origData := make([]byte, len(cipherText))

	switch c.mode {
	case ModeCBC:
		enc := cipher.NewCBCDecrypter(c.block, c.iv)
		enc.CryptBlocks(origData, cipherText)
	case ModeECB:
		enc := ecb.NewECBDecrypter(c.block)
		enc.CryptBlocks(origData, cipherText)
	case ModeCFB:
		enc := cipher.NewCFBDecrypter(c.block, c.iv)
		enc.XORKeyStream(origData, cipherText)
	case ModeOFB:
		enc := cipher.NewOFB(c.block, c.iv)
		enc.XORKeyStream(origData, cipherText)
	default:
		return "", errors.New("invalid cipher mode")
	}
	origData, err = c.padding.UnPad(origData)
	if err != nil {
		return "", err
	}
	return strings.BytesToString(origData), nil
}
