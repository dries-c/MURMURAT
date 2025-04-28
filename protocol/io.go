package protocol

type IO interface {
	Uint8(x *uint8)
	BEUint16(x *uint16)
	BEUint32(x *uint32)
	Bytes(x *[]byte, n int)
}
