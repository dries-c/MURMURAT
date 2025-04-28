package message

import "MURMURAT/protocol"

type DataMessage struct {
	nonce       []byte
	timestamp   uint32
	data        []byte
	publicKeyId []byte
	signature   []byte
}

func (x *DataMessage) ID() uint8 {
	return IDData
}

func (x *DataMessage) Marshal(r protocol.IO) error {
	var length uint16
	r.BEUint16(&length)
	r.Bytes(&x.nonce, 1)
	r.BEUint32(&x.timestamp)
	r.Bytes(&x.data, int(length-1-4-4-512)) // Stupid implementation
	r.Bytes(&x.publicKeyId, 4)
	r.Bytes(&x.signature, 512)
	return nil
}
