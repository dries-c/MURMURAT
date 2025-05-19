package mitm

import (
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"net"
)

type UDP struct {
	conn    *net.UDPConn
	snooper func(udp *UDP, buf []byte, addr net.Addr, incoming bool) bool
}

func NewUDP(src net.IP, port int) (*UDP, error) {
	addr := &net.UDPAddr{
		IP:   src,
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDP{conn: conn}, nil
}

func (udp *UDP) Read(buf []byte) (int, net.Addr, error) {
	n, addr, err := udp.conn.ReadFromUDP(buf)

	if err != nil {
		return 0, nil, err
	}

	if udp.snooper != nil {
		if !udp.snooper(udp, buf[:n], addr, true) {
			return 0, nil, nil
		}
	}

	return n, addr, nil
}

func (udp *UDP) Write(buf []byte, addr net.Addr) (int, error) {
	n, err := udp.conn.WriteToUDP(buf, addr.(*net.UDPAddr))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (udp *UDP) registerSnooper(snooper func(udp *UDP, buf []byte, addr net.Addr, incoming bool) bool) {
	udp.snooper = snooper
}

func (udp *UDP) WritePacket(msg message.Message, addr net.Addr, ignoreSnooper bool) (int, error) {
	buf := bytes.NewBuffer(nil)
	writer := protocol.NewWriter(buf)

	header := &message.Header{
		ID: msg.ID(),
	}

	if err := header.Marshal(writer); err != nil {
		return 0, err
	}

	if err := msg.Marshal(writer); err != nil {
		return 0, err
	}

	if !ignoreSnooper && udp.snooper != nil {
		if !udp.snooper(udp, buf.Bytes(), addr, false) {
			return 0, nil
		}
	}

	return udp.Write(buf.Bytes(), addr)
}
