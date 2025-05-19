package mitm

import (
	"fmt"
	"net"
)

type Client struct {
	base   *Base
	target net.Addr
}

func NewClient(target net.Addr, srcPort int, onSessionInitialized func(session *Session) error) *Client {
	base := newBase(srcPort)
	base.setOnSessionInitialized(func(session *Session, client bool) error {
		if !client {
			return nil
		}

		session.SetOnDataReceived(func(data []byte) error {
			fmt.Println("Client received data from server:", string(data))
			return nil
		})

		return onSessionInitialized(session)
	})

	return &Client{
		base:   base,
		target: target,
	}
}

func (c *Client) RegisterSnooper(snooper func(udp *UDP, buf []byte, addr net.Addr, incoming bool) bool) {
	c.base.udp.registerSnooper(snooper)
}

func (c *Client) Start() {
	c.base.getSession(c.target, true)
	c.base.start()
}
