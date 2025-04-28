package message

import "MURMURAT/protocol"

const (
	IDDH = iota
	IDHello
	IDData
)

type Message interface {
	ID() uint8
	Marshal(io protocol.IO) error
}

type Header struct {
	ID uint8
}

func (x *Header) Marshal(r protocol.IO) error {
	r.Uint8(&x.ID)
	return nil
}
