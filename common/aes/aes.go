package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var key = [32]byte{
	0x7A, 0x2C, 0x74, 0x6D, 0xB5, 0x4F, 0xF7, 0xAF, 0x4A, 0x18, 0x8D, 0x94, 0x7A, 0xE4, 0x71, 0x01,
	0x44, 0x19, 0xE6, 0x83, 0x68, 0x46, 0x86, 0xDB, 0xBE, 0x6D, 0xD9, 0x9C, 0x8C, 0x3C, 0x08, 0x40,
}

// Aes wrap
type Aes struct {
	cb cipher.Block
}

// New create a new instance
func New(key []byte) *Aes {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	return &Aes{block}
}

// Encrypt encrypt buf
func (c *Aes) Encrypt(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("src len is zero")
	}

	// padding src
	padsize := aes.BlockSize - len(src)%aes.BlockSize
	bufpad := bytes.Repeat([]byte{0}, padsize)
	dst := append(src, bufpad...)

	// encrypt
	for i := 0; i < len(dst); i += aes.BlockSize {
		c.cb.Encrypt(dst[i:], dst[i:])
	}

	// fill padsize
	dst = append(dst, uint8(padsize))
	return dst, nil
}

// Decrypt decrypt buf, return buf decrypted
func (c *Aes) Decrypt(src []byte) ([]byte, error) {
	if len(src)%aes.BlockSize != 1 {
		return nil, errors.New("src len invalid")
	}

	// extract padsize
	padsize := int(src[len(src)-1])
	dst := make([]byte, len(src)-1)

	// decrypt
	for i := 0; i < len(dst); i += aes.BlockSize {
		c.cb.Decrypt(dst[i:], src[i:])
	}
	dst = dst[:len(dst)-padsize]
	return dst, nil
}

var defaultAES = New(key[:])

// Encrypt encrypt buf
func Encrypt(src []byte) ([]byte, error) {
	return defaultAES.Encrypt(src)
}

// Decrypt decrypt buf, return buf decrypted
func Decrypt(src []byte) ([]byte, error) {
	return defaultAES.Decrypt(src)
}
