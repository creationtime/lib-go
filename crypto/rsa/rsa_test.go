package rsa

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	privKey = `MIICeQIBADANBgkqhkiG9w0BAQEFAASCAmMwggJfAgEAAoGBAMheaSuMWdPgeL2Lras8qPyINGYIOJlYQR/Az4411YSI7ghUZVI9JBVO8Mp3+qroxwLPqXSWx9TZORvxsZy1dSvMtYZihy1s5aUhI7jCfkJYu3KwRVirLQeviA7Fvh0c/2XGFPoTmrsGrPawL5pvmv1zjbTGCOxTWb0zQvqwebIfAgMBAAECgYEAxNAn0gTcv0gAkX7AKjE9dEB957MvlUChR4Vm2rN6deLinP/5PlycMuoFj3tml7ZqtRIxyznINATjGdXAtsNuwMfjBMWtFfit48xiQxZGweZ2X3fA9pQ5MMnHoIJbhdSpBxJ/CCFFP55zUDKm6Ydfk9rJutPKTdQKv6m+VCVhfgECQQDZ+mj88yxD+NMwp1ykPKxkwYkABN7BXT0ha3r0VhNorN88FqVMAxhdOgSRXQTIPyCtgj9KPdeWJ44kKU++g5GBAkEA61GsIL+5+AT5F+9cpiB3ckj2kpt8rzP+XlNBuMg8KSoE8j1odqLKaCejpq4kDW4QKOyyPAaA14kqKmWTnFXTnwJBAM/BAlWsc5EpVCg4K20BwxGZADl7atADTONQbHT6oS8QLQg5UTx8arlYNchSTt+Ig128GRRqktKzSp+enDmpboECQQCfwu1HtqM9nbK360xNhVFTB/JPirzV/ki+JWxDVb5yfBKrm8FmehNNL0xOB4B2lbjm7/v6ALhMnNVBv4C97Q8lAkEAu3gFnapIBkLogVccesvdxI9T+hJtws9ttAA5MUltkdhnhAr52RKHf8Hedi3JIdCDN1RrqowbDrJW5dW1V9MjMw==`
	pubKey  = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIXmkrjFnT4Hi9i62rPKj8iDRmCDiZWEEfwM+ONdWEiO4IVGVSPSQVTvDKd/qq6McCz6l0lsfU2Tkb8bGctXUrzLWGYoctbOWlISO4wn5CWLtysEVYqy0Hr4gOxb4dHP9lxhT6E5q7Bqz2sC+ab5r9c420xgjsU1m9M0L6sHmyHwIDAQAB`

	r *Rsa

	text        = "加餐聊顿觉，密意未全消，解榻情何限，密疏坐寂寥。"
	cryptedText = "UotCVgUVuMR4+nYVF7/tF0/i1DcfZyYn43FTzqUDzcJZyvKhuH7bNEBuoYt29P16/MsVXJQzhVBmJYGvvaVGih0FULf2KzzuaLfg2CO+CxsDFvX5y3t8VS/qsZcx0gJIBy/WbH52WCq5jevg3BArVClZEeTWz8yrEQvFNkClc3A="
	sign        = "ADxvkSqKPrrjloyiwRUna4zRcY9AmfaAyWla1K7SrS20XT3V4nTxP9z4JD5WGDq0EGley/9loHvN1nW9E72CJ7ZaCTongImkzO3k+GzxhmIjCCF6y1/H3DWjPpXPL81Qtzg+S4B4h89kEcuTpHK9uzx3OAuB2TYzmiKhRGcZg7s="
)

func init() {
	var err error
	if r, err = New(pubKey, privKey); err != nil {
		log.Fatal(err)
	}
}

func TestGenerateKey(t *testing.T) {
	err := GenerateKey(1024)
	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	r, err := New(pubKey, privKey)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestRsa_Encrypt(t *testing.T) {
	result, err := r.Encrypt(text)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func BenchmarkRsa_Encrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := r.Encrypt(text)
		assert.NoError(b, err)
	}
}

func TestRsa_Decrypt(t *testing.T) {
	result, err := r.Decrypt(cryptedText)
	assert.NoError(t, err)
	assert.EqualValues(t, text, result)
}

func BenchmarkRsa_Decrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := r.Decrypt(cryptedText)
		assert.NoError(b, err)
	}
}

func TestRsa_Sign(t *testing.T) {
	sign, err := r.Sign(text)
	assert.NoError(t, err)
	assert.NotEmpty(t, sign)
}

func BenchmarkRsa_Sign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := r.Sign(text)
		assert.NoError(b, err)
	}
}

func TestRsa_Verify(t *testing.T) {
	err := r.Verify(text, sign)
	assert.NoError(t, err)
}

func BenchmarkRsa_Verify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := r.Verify(text, sign)
		assert.NoError(b, err)
	}
}
