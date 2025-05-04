package handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/binary"
	"math/big"
)

type SignatureVerifier struct {
	publicKey *rsa.PublicKey
}

type SignatureCreator struct {
	privateKey *rsa.PrivateKey
}

func NewSignatureCreator() *SignatureCreator {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic("failed to generate private key: " + err.Error())
	}

	privateKey.E = 65537
	return &SignatureCreator{
		privateKey: privateKey,
	}
}

func NewSignatureVerifier(publicKey *rsa.PublicKey) *SignatureVerifier {
	return &SignatureVerifier{
		publicKey: publicKey,
	}
}

func NewSignatureVerifierFromBytes(modulusBytes []byte) *SignatureVerifier {
	modulus := new(big.Int).SetBytes(modulusBytes)
	publicKey := &rsa.PublicKey{
		N: modulus,
		E: 65537,
	}

	return &SignatureVerifier{
		publicKey: publicKey,
	}
}

func hashData(data []byte) ([]byte, error) {
	hash := sha3.New256()

	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func (h *SignatureCreator) GetPublicKey() *rsa.PublicKey {
	return &h.privateKey.PublicKey
}

func (h *SignatureCreator) GetPublicKeyId() uint32 {
	hash := sha256.Sum256(h.privateKey.PublicKey.N.Bytes())
	return binary.BigEndian.Uint32(hash[:4])
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
