//go:build !armbe && !arm64be && !ppc64 && !mips && !mips64 && !mips64p32 && !ppc && !sparc && !sparc64 && !s390 && !s390x

package protocol

import "encoding/binary"

func (w *Writer) BEUint16(x *uint16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, *x)
	_, _ = w.w.Write(data)
}

func (w *Writer) BEUint32(x *uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, *x)
	_, _ = w.w.Write(data)
}
