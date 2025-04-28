package handler

import (
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"fmt"
)

type PacketHander struct {
	pool message.Pool
}

func NewPacketHandler() *PacketHander {
	return &PacketHander{
		pool: message.NewPool(),
	}
}

func (h *PacketHander) Handle(data []byte) error {
	buf := bytes.NewBuffer(data)
	reader := protocol.NewReader(buf)

	header := &message.Header{}
	if err := header.Marshal(reader); err != nil {
		return fmt.Errorf("read message header: %w", err)
	}

	msg, err := h.pool.Get(header.ID)
	if err != nil {
		return fmt.Errorf("get message: %w", err)
	}

	if err := msg.Marshal(reader); err != nil {
		return fmt.Errorf("read message data: %w", err)
	}

	fmt.Printf("Received message with ID: %d\n", msg.ID())
	// print message data
	fmt.Printf("Message data: %v\n", msg)
	return nil
}
