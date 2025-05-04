package protocol

import (
	"crypto/rsa"
	"io"
	"math/big"
)

type Reader struct {
	r interface {
		io.Reader
		io.ByteReader
	}
}

func NewReader(r interface {
	io.Reader
	io.ByteReader
}) *Reader {
	return &Reader{r: r}
}

func (r *Reader) Uint8(x *uint8) {
	var err error
	*x, err = r.r.ReadByte()
	if err != nil {
		r.panic(err)
	}
}

func (r *Reader) Bytes(x *[]byte, n int) {
	b := make([]byte, n)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = b
}

func (r *Reader) RSAPublicKey(key **rsa.PublicKey) {
	b := make([]byte, 512)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}

	*key = &rsa.PublicKey{
		N: new(big.Int).SetBytes(b),
		E: 65537,
	}
}

func (r *Reader) panic(err error) {
	panic(err)
}
