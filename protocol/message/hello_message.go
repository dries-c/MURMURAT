package message

import "MURMURAT/protocol"

type HelloMessage struct {
	PublicKeyId  []byte
	RsaPublicKey []byte
}

func (x *HelloMessage) ID() uint8 {
	return IDHello
}

func (x *HelloMessage) Marshal(r protocol.IO) error {
	r.Bytes(&x.PublicKeyId, 4)
	r.Bytes(&x.RsaPublicKey, 512)
	return nil
}
