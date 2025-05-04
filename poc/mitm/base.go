package mitm

import (
	"MURMURAT/protocol/message"
	"log"
	"net"
)

type Base struct {
	sessions             map[string]*Session
	udp                  *UDP
	onSessionCreated     func(session *Session, client bool)
	onSessionInitialized func(session *Session, client bool) error
}

func NewBase(srcPort int) *Base {
	udp, err := NewUDP(net.IPv4(0, 0, 0, 0), srcPort)
	if err != nil {
		log.Fatalf("Error creating UDP connection: %v", err)
	}

	return &Base{
		sessions: make(map[string]*Session),
		udp:      udp,
	}
}

func (s *Base) Start() {
	bytes := make([]byte, 1024)

	for {
		read, addr, err := s.udp.Read(bytes)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		session := s.getSession(addr, false)
		err = session.PacketHandler.Handle(bytes[:read])
		if err != nil {
			log.Printf("Error handling packet: %v", err)
		}
	}
}

func (s *Base) getSession(addr net.Addr, client bool) *Session {
	session, exists := s.sessions[addr.String()]
	if exists {
		return session
	}

	newSession := NewSession(addr, s.udp)

	if client {
		err := newSession.SendDHMessage()
		if err != nil {
			return nil
		}

		s.registerClientListeners(newSession)
	} else {
		s.registerServerListeners(newSession)
	}

	if s.onSessionCreated != nil {
		s.onSessionCreated(newSession, client)
	}

	s.sessions[addr.String()] = newSession
	return newSession
}

func (s *Base) SetOnNewSessionCreated(callback func(session *Session, client bool)) {
	s.onSessionCreated = callback
}

func (s *Base) SetOnSessionInitialized(callback func(session *Session, client bool) error) {
	s.onSessionInitialized = callback
}

func (s *Base) registerServerListeners(session *Session) {
	packetHandler := session.PacketHandler
	packetHandler.RegisterListener(message.IDDH, func(msg message.Message) error {
		return session.SendDHMessage()
	})
	packetHandler.RegisterListener(message.IDHello, func(msg message.Message) error {
		err := session.SendHelloMessage()
		if err != nil {
			return err
		}

		if s.onSessionInitialized != nil {
			return s.onSessionInitialized(session, false)
		}

		return nil
	})
}

func (s *Base) registerClientListeners(session *Session) {
	packetHandler := session.PacketHandler
	packetHandler.RegisterListener(message.IDDH, func(msg message.Message) error {
		return session.SendHelloMessage()
	})
	packetHandler.RegisterListener(message.IDHello, func(msg message.Message) error {
		if s.onSessionInitialized != nil {
			return s.onSessionInitialized(session, true)
		}

		return nil
	})
}
