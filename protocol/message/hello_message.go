package message

import "MURMURAT/protocol"

type HelloMessage struct {
	publicKeyId  []byte
	rsaPublicKey []byte
}

func (x *HelloMessage) ID() uint8 {
	return IDHello
}

func (x *HelloMessage) Marshal(r protocol.IO) error {
	r.Bytes(&x.publicKeyId, 4)
	r.Bytes(&x.rsaPublicKey, 512)
	return nil
}
