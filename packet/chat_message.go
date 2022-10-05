package packet

type ChatMessageServer struct {
	Message string
}

func (m ChatMessageServer) ID() int32 {
	return 0x03
}

func (m *ChatMessageServer) Decode(r *Reader) error {
	return r.String(&m.Message)
}

func (m ChatMessageServer) Encode(w Writer) error {
	return w.String(m.Message)
}
