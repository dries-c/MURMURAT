package mitm

import (
	"MURMURAT/handler"
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"fmt"
	"net"
	"time"
)

type Session struct {
	publicKeyId uint32
	addr        net.Addr

	udp *UDP

	PacketHandler *handler.PacketHandler
	dh            *protocol.DiffieHellman

	signatureVerifier *handler.SignatureVerifier
	signatureCreator  *handler.SignatureCreator

	encryptionHandler *handler.EncryptionHandler
	onDataReceived    func(data []byte) error
}

func NewSession(addr net.Addr, udp *UDP) *Session {
	s := &Session{
		addr: addr,

		udp: udp,

		PacketHandler: handler.NewPacketHandler(),
		dh:            protocol.NewDiffieHellman(),

		signatureCreator: handler.NewSignatureCreator(),
	}

	s.registerBaseListener()

	return s
}

func (s *Session) registerBaseListener() {
	s.PacketHandler.RegisterValidator(func(msg message.Message) error {
		switch msg.(type) {
		case *message.DataMessage:
			if s.signatureVerifier == nil {
				return fmt.Errorf("signature verifier not initialized")
			}

			dataMessage, ok := msg.(*message.DataMessage)
			if !ok {
				return fmt.Errorf("invalid message type")
			}

			t := time.Unix(int64(dataMessage.Timestamp), 0)
			if t.After(time.Now().Add(1 * time.Minute)) {
				return fmt.Errorf("timestamp is older than 1 minute")
			}

			if err := s.signatureVerifier.Verify(dataMessage.Data, dataMessage.Signature); err != nil {
				return fmt.Errorf("signature verification failed: %w", err)
			}
			break
		}
		return nil
	})
	s.PacketHandler.RegisterListener(message.IDDH, func(msg message.Message) error {
		dhMessage, ok := msg.(*message.DHMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		err := s.dh.ComputeSharedKey(dhMessage.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to compute shared key: %w", err)
		}

		s.encryptionHandler = handler.NewEncryptionHandler(s.dh.SessionKey)
		return nil
	})

	s.PacketHandler.RegisterListener(message.IDHello, func(msg message.Message) error {
		helloMessage, ok := msg.(*message.HelloMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		s.publicKeyId = helloMessage.PublicKeyId
		s.signatureVerifier = handler.NewSignatureVerifier(helloMessage.RsaPublicKey)
		return nil
	})

	s.PacketHandler.RegisterListener(message.IDData, func(msg message.Message) error {
		dataMessage, ok := msg.(*message.DataMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		decrypt, err := s.encryptionHandler.Decrypt(dataMessage.Data, dataMessage.Nonce)
		if err != nil {
			return err
		}

		if s.onDataReceived != nil {
			return s.onDataReceived(decrypt)
		}
		return nil
	})
}

func (s *Session) SetOnDataReceived(callback func(data []byte) error) {
	s.onDataReceived = callback
}

func (s *Session) SendMessage(msg message.Message) error {
	_, err := s.udp.WritePacket(msg, s.addr, false)
	return err
}

func (s *Session) SendDHMessage() error {
	return s.SendMessage(&message.DHMessage{
		PublicKey: s.dh.PublicKey.Bytes(),
	})
}

func (s *Session) SendHelloMessage() error {
	return s.SendMessage(&message.HelloMessage{
		PublicKeyId:  s.signatureCreator.GetPublicKeyId(),
		RsaPublicKey: s.signatureCreator.GetPublicKey(),
	})
}

func (s *Session) SendDataMessage(data []byte) error {
	nonce, err := s.encryptionHandler.GenerateNonce()
	if err != nil {
		return err
	}

	encryptedData, err := s.encryptionHandler.Encrypt(data, nonce)
	if err != nil {
		return err
	}

	signature, err := s.signatureCreator.Sign(encryptedData)
	if err != nil {
		return err
	}

	return s.SendMessage(&message.DataMessage{
		Nonce:       nonce,
		Timestamp:   uint32(time.Now().Unix()),
		Data:        encryptedData,
		PublicKeyId: s.signatureCreator.GetPublicKeyId(),
		Signature:   signature,
	})
}
