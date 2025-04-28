package handler

import (
	"crypto/aes"
	"crypto/cipher"
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

func (e *EncryptionHandler) Encrypt(plaintext []byte, nonce []byte) ([]byte, error) {
	return e.xORKeyStream(plaintext, nonce)
}

func (e *EncryptionHandler) Decrypt(ciphertext []byte, nonce []byte) ([]byte, error) {
	return e.xORKeyStream(ciphertext, nonce)
}

func (e *EncryptionHandler) xORKeyStream(streamText []byte, nonce []byte) ([]byte, error) {
	// IV: 15 zero bytes + 1-byte nonce at the end
	iv := make([]byte, 16)
	iv[15] = nonce[0]

	stream := cipher.NewCTR(e.block, iv)
	text := make([]byte, len(streamText))
	stream.XORKeyStream(text, streamText)

	return text, nil
}

func XORBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("lengths of a and b must be equal")
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
