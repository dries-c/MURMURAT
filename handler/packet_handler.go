package handler

import (
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"fmt"
)

type PacketHander struct {
	pool     message.Pool
	listener map[uint8][]func(message.Message) error
}

func NewPacketHandler() *PacketHander {
	return &PacketHander{
		pool:     message.NewPool(),
		listener: make(map[uint8][]func(message.Message) error),
	}
}

func (h *PacketHander) RegisterListener(id uint8, listener func(message.Message) error) {
	if _, exists := h.listener[id]; !exists {
		h.listener[id] = []func(message.Message) error{}
	}
	h.listener[id] = append(h.listener[id], listener)
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
	fmt.Printf("Message data: %v\n", msg)

	if listeners, exists := h.listener[msg.ID()]; exists {
		for _, listener := range listeners {
			if err := listener(msg); err != nil {
				return fmt.Errorf("listener error: %w", err)
			}
		}
	}

	return nil
}
