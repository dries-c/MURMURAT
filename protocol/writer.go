package protocol

import (
	"crypto/rsa"
	"io"
)

type Writer struct {
	w interface {
		io.Writer
		io.ByteWriter
	}
}

func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Uint8(x *uint8) {
	if err := w.w.WriteByte(*x); err != nil {
		w.panic(err)
	}
}

func (w *Writer) Bytes(x *[]byte, n int) {
	if _, err := w.w.Write((*x)[:n]); err != nil {
		w.panic(err)
	}
}

func (w *Writer) RSAPublicKey(key **rsa.PublicKey) {
	paddedKey := make([]byte, 512)
	copy(paddedKey, (*key).N.Bytes())
	w.Bytes(&paddedKey, 512)
}

func (w *Writer) panic(err error) {
	panic(err)
}
