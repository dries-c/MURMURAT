package message

import "MURMURAT/protocol"

type DataMessage struct {
	Nonce       uint8
	Timestamp   uint32
	Data        []byte
	PublicKeyId []byte
	Signature   []byte
}

func (x *DataMessage) ID() uint8 {
	return IDData
}

func (x *DataMessage) Marshal(r protocol.IO) error {
	var length uint16
	r.BEUint16(&length)
	r.Uint8(&x.Nonce)
	r.BEUint32(&x.Timestamp)
	r.Bytes(&x.Data, int(length-1-4-4-512)) // Stupid protocol spec
	r.Bytes(&x.PublicKeyId, 4)
	r.Bytes(&x.Signature, 512)
	return nil
}
