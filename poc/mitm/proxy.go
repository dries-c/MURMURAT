package mitm

import (
	"MURMURAT/protocol/message"
	"fmt"
	"net"
)

type Proxy struct {
	server *Base
	target net.Addr
}

func NewProxy(target net.Addr, srcPort int) *Proxy {
	server := newBase(srcPort)
	server.setOnNewSessionCreated(func(session *Session, client bool) {
		if client {
			return
		}

		isConnected := false
		var cache [][]byte

		clientSession := server.getSession(target, true)
		server.setOnSessionInitialized(func(session *Session, client bool) error {
			if !client || isConnected || session != clientSession {
				return nil
			}

			isConnected = true
			for _, data := range cache {
				if err := session.SendDataMessage(data); err != nil {
					return err
				}
			}

			return nil
		})
		session.PacketHandler.RegisterListener(message.IDDH, func(msg message.Message) error {
			return clientSession.SendDHMessage()
		})
		session.PacketHandler.RegisterListener(message.IDHello, func(msg message.Message) error {
			return clientSession.SendHelloMessage()
		})
		session.SetOnDataReceived(func(data []byte) error {
			fmt.Println("Proxy received data from client:", string(data))

			if !isConnected {
				cache = append(cache, data)
				return nil
			}

			return clientSession.SendDataMessage(data)
		})
	})

	return &Proxy{
		server: server,
		target: target,
	}
}

func (c *Proxy) Start() {
	c.server.start()
}
