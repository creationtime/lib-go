package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/creationtime/lib-go/util/strings"
)

type Rsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	encoding   *base64.Encoding
	h          crypto.Hash
}

func GenerateKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	privateStream, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	fmt.Println(base64.StdEncoding.EncodeToString(privateStream))

	publicKey := privateKey.PublicKey
	publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	fmt.Println(base64.StdEncoding.EncodeToString(publicStream))
	return nil
}

func New(pubKey string, privKey string) (*Rsa, error) {
	b, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	b, err = base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	priv, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		return nil, err
	}
	r := &Rsa{
		publicKey:  pub.(*rsa.PublicKey),
		privateKey: priv.(*rsa.PrivateKey),
		encoding:   base64.StdEncoding,
		h:          crypto.SHA256,
	}
	return r, nil
}

func (r *Rsa) Encrypt(src string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split(strings.StringToBytes(src), partLen)

	var buf bytes.Buffer
	for _, chunk := range chunks {
		b, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buf.Write(b)
	}
	return r.encoding.EncodeToString(buf.Bytes()), nil
}

func (r *Rsa) Decrypt(src string) (string, error) {
	raw, err := r.encoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	partLen := r.publicKey.N.BitLen() / 8
	chunks := split(raw, partLen)

	var buf bytes.Buffer
	for _, chunk := range chunks {
		b, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buf.Write(b)
	}
	return buf.String(), err
}

func (r *Rsa) hash(src string) []byte {
	h := r.h.New()
	h.Write(strings.StringToBytes(src))
	return h.Sum(nil)
}

func (r *Rsa) Sign(src string) (string, error) {
	hashed := r.hash(src)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, r.h, hashed)
	if err != nil {
		return "", err
	}
	return r.encoding.EncodeToString(sign), err
}

func (r *Rsa) Verify(src string, sign string) error {
	hashed := r.hash(src)

	decodedSign, err := r.encoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.publicKey, r.h, hashed, decodedSign)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
