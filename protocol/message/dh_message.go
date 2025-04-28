package message

import "MURMURAT/protocol"

type DHMessage struct {
	publicKey []byte
}

func (x *DHMessage) ID() uint8 {
	return IDDH
}

func (x *DHMessage) Marshal(r protocol.IO) error {
	r.Bytes(&x.publicKey, 256)
	return nil
}
