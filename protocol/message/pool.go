package message

type Pool map[uint8]func() Message

func NewPool() Pool {
	pool := make(Pool)

	pool.Register(IDDH, func() Message { return &DHMessage{} })
	pool.Register(IDHello, func() Message { return &HelloMessage{} })
	pool.Register(IDData, func() Message { return &DataMessage{} })

	return pool
}

func (p Pool) Get(id uint8) (Message, error) {
	if constructor, ok := p[id]; ok {
		return constructor(), nil
	}
	return nil, ErrUnknownPacket
}

func (p Pool) Register(id uint8, constructor func() Message) {
	if _, exists := p[id]; exists {
		panic("message already registered")
	}
	p[id] = constructor
}
