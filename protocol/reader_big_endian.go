//go:build armbe || arm64be || ppc64 || mips || mips64 || mips64p32 || ppc || sparc || sparc64 || s390 || s390x

package protocol

import "unsafe"

func (r *Reader) BEUint16(x *uint16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = *(*uint16)(unsafe.Pointer(&b[0]))
}

func (r *Reader) BEUint32(x *uint32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = *(*uint32)(unsafe.Pointer(&b[0]))
}
