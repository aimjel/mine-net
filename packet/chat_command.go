package packet

type ChatCommandServer struct {
	Command              string
	Timestamp            int64
	Salt                 int64
	ArgumentSignatures   []Argument
	AcknowledgedMessages []int64
}

type Argument struct {
	Name      string
	Signature []byte
}

func (m ChatCommandServer) ID() int32 {
	return 0x04
}

func (m *ChatCommandServer) Decode(r *Reader) error {
	r.String(&m.Command)
	r.Int64(&m.Timestamp)
	r.Int64(&m.Salt)
	var length int32
	r.VarInt(&length)
	for i := int32(0); i < length; i++ {
		var name string
		var sig = make([]byte, 256)
		r.String(&name)
		r.FixedByteArray(&sig)
		m.ArgumentSignatures = append(m.ArgumentSignatures, Argument{
			Name:      name,
			Signature: sig,
		})
	}
	var count int32
	r.VarInt(&count)
	m.AcknowledgedMessages = make([]int64, count)
	for i := int32(0); i < count; i++ {
		r.Int64(&m.AcknowledgedMessages[i])
	}
	return nil
}

func (m ChatCommandServer) Encode(w *Writer) error {
	w.String(m.Command)
	w.Int64(m.Timestamp)
	w.Int64(m.Salt)
	w.VarInt(int32(len(m.ArgumentSignatures)))
	for _, a := range m.ArgumentSignatures {
		w.String(a.Name)
		w.FixedByteArray(a.Signature)
	}
	w.VarInt(int32(len(m.AcknowledgedMessages)))
	for _, a := range m.AcknowledgedMessages {
		w.Int64(a)
	}
	return nil
}
