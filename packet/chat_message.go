package packet

type ChatMessageServer struct {
	Message string

	//TODO add the rest of the fields
	//https://wiki.vg/Protocol#Chat_Message
}

func (m ChatMessageServer) ID() int32 {
	return 0x05
}

func (m *ChatMessageServer) Decode(r *Reader) error {
	return r.String(&m.Message)
}

func (m ChatMessageServer) Encode(w Writer) error {
	return w.String(m.Message)
}
