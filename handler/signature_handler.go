package handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha3"
	"math/big"
)

type SignatureVerifier struct {
	publicKey *rsa.PublicKey
}

type SignatureCreator struct {
	privateKey *rsa.PrivateKey
}

func NewSignatureCreator(privateKey *rsa.PrivateKey) *SignatureCreator {
	return &SignatureCreator{
		privateKey: privateKey,
	}
}

func NewSignatureVerifier(publicKey *rsa.PublicKey) *SignatureVerifier {
	return &SignatureVerifier{
		publicKey: publicKey,
	}
}

func RSAPublicKeyFromBytes(modulusBytes []byte, exponent int) *rsa.PublicKey {
	modulus := new(big.Int).SetBytes(modulusBytes)
	return &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}
}

func hashData(data []byte) ([]byte, error) {
	hash := sha3.New256()

	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	hashed := hash.Sum(nil)
	return hashed, nil
}

func (h *SignatureCreator) Sign(data []byte) ([]byte, error) {
	hash, err := hashData(data)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, h.privateKey, crypto.SHA3_256, hash)
}

func (h *SignatureVerifier) Verify(data []byte, sig []byte) error {
	hash, err := hashData(data)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(h.publicKey, crypto.SHA3_256, hash, sig)
}
