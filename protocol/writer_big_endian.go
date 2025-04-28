//go:build armbe || arm64be || ppc64 || mips || mips64 || mips64p32 || ppc || sparc || sparc64 || s390 || s390x

package protocol

import (
	"unsafe"
)

func (w *Writer) BEUint16(x *uint16) {
	data := *(*[2]byte)(unsafe.Pointer(x))
	_, _ = w.w.Write(data[:])
}

func (w *Writer) BEUint32(x *uint32) {
	data := *(*[4]byte)(unsafe.Pointer(x))
	_, _ = w.w.Write(data[:])
}
