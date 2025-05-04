package message

import (
	"MURMURAT/protocol"
	"crypto/rsa"
)

type HelloMessage struct {
	PublicKeyId  uint32
	RsaPublicKey *rsa.PublicKey
}

func (x *HelloMessage) ID() uint8 {
	return IDHello
}

func (x *HelloMessage) Marshal(r protocol.IO) error {
	r.BEUint32(&x.PublicKeyId)
	r.RSAPublicKey(&x.RsaPublicKey)
	return nil
}
