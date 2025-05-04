package mitm

import "fmt"

type Server struct {
	base *Base
}

func NewServer(srcPort int) *Server {
	base := NewBase(srcPort)
	base.SetOnNewSessionCreated(func(session *Session, client bool) {
		if client {
			return
		}

		session.SetOnDataReceived(func(data []byte) error {
			fmt.Println("Server received data from client:", string(data))
			return nil
		})
	})

	return &Server{
		base: base,
	}
}

func (c *Server) Start() {
	c.base.Start()
}
