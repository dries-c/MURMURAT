package mitm

import (
	"fmt"
)

type Server struct {
	base *Base
}

func NewServer(srcPort int) *Server {
	base := newBase(srcPort)
	base.setOnNewSessionCreated(func(session *Session, client bool) {
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

func (s *Server) SetOnSessionInitialized(onSessionCreated func(session *Session)) {
	s.base.onSessionInitialized = func(session *Session, client bool) error {
		if client {
			return nil
		}

		if onSessionCreated != nil {
			onSessionCreated(session)
		}
		return nil
	}
}

func (c *Server) Start() {
	c.base.start()
}
