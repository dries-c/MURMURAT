package message

import "MURMURAT/protocol"

type DHMessage struct {
	PublicKey []byte
}

func (x *DHMessage) ID() uint8 {
	return IDDH
}

func (x *DHMessage) Marshal(r protocol.IO) error {
	r.Bytes(&x.PublicKey, 256)
	return nil
}
