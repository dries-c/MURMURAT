package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type EncryptionHandler struct {
	block cipher.Block
}

func NewEncryptionHandler(key []byte) *EncryptionHandler {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return &EncryptionHandler{
		block: block,
	}
}

func (e *EncryptionHandler) Encrypt(plaintext []byte, nonce byte) ([]byte, error) {
	return e.xORKeyStream(plaintext, nonce)
}

func (e *EncryptionHandler) Decrypt(ciphertext []byte, nonce byte) ([]byte, error) {
	return e.xORKeyStream(ciphertext, nonce)
}

func (e *EncryptionHandler) xORKeyStream(streamText []byte, nonce byte) ([]byte, error) {
	// IV: 15 zero bytes + 1-byte nonce at the end
	iv := make([]byte, 16)
	iv[15] = nonce

	stream := cipher.NewCTR(e.block, iv)
	text := make([]byte, len(streamText))
	stream.XORKeyStream(text, streamText)

	return text, nil
}

func (e *EncryptionHandler) GenerateNonce() (byte, error) {
	var b [1]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func XORBytes(a, b []byte) []byte {
	minLength := len(a)
	if len(b) < minLength {
		minLength = len(b)
	}

	result := make([]byte, minLength)
	for i := 0; i < minLength; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}
