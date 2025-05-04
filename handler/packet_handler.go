package handler

import (
	"MURMURAT/protocol"
	"MURMURAT/protocol/message"
	"bytes"
	"fmt"
)

type PacketHandler struct {
	pool      message.Pool
	listener  map[uint8][]func(message.Message) error
	validator func(message.Message) error
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		pool:     message.NewPool(),
		listener: make(map[uint8][]func(message.Message) error),
	}
}

func (h *PacketHandler) RegisterValidator(validator func(message.Message) error) {
	if h.validator != nil {
		panic("validator already registered")
	}

	h.validator = validator
}

func (h *PacketHandler) RegisterListener(id uint8, listener func(message.Message) error) {
	if _, exists := h.listener[id]; !exists {
		h.listener[id] = []func(message.Message) error{}
	}
	h.listener[id] = append(h.listener[id], listener)
}

func (h *PacketHandler) Handle(data []byte) error {
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

	if h.validator != nil {
		if err := h.validator(msg); err != nil {
			return fmt.Errorf("invalid message: %w", err)
		}
	}

	if listeners, exists := h.listener[msg.ID()]; exists {
		for _, listener := range listeners {
			if err := listener(msg); err != nil {
				return fmt.Errorf("listener error: %w", err)
			}
		}
	}

	return nil
}
