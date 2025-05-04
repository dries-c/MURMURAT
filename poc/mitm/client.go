package mitm

import (
	"net"
)

type Client struct {
	base   *Base
	target net.Addr
}

func NewClient(target net.Addr, srcPort int, onSessionInitialized func(session *Session) error) *Client {
	base := NewBase(srcPort)
	base.SetOnSessionInitialized(func(session *Session, client bool) error {
		if !client {
			return nil
		}

		return onSessionInitialized(session)
	})

	return &Client{
		base:   base,
		target: target,
	}
}

func (c *Client) Start() {
	c.base.getSession(c.target, true)
	c.base.Start()
}
