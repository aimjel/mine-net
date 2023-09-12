package packet

type ChatCommandServer struct {
	Command   string
	Timestamp int64
	Salt      int64

	//TODO add the rest of the fields
	//https://wiki.vg/Protocol#Chat_Command
}

func (m ChatCommandServer) ID() int32 {
	return 0x04
}

func (m *ChatCommandServer) Decode(r *Reader) error {
	r.String(&m.Command)
	r.Int64(&m.Timestamp)
	return r.Int64(&m.Salt)
}

func (m ChatCommandServer) Encode(w Writer) error {
	w.String(m.Command)
	w.Int64(m.Timestamp)
	return w.Int64(m.Salt)
}
