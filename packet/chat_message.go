package packet

import "github.com/aimjel/minenet/protocol/encoding"

type ChatMessageServer struct {
	Message              string
	Timestamp            int64
	Salt                 int64
	Signature            []byte
	AcknowledgedMessages []int64
}

func (m ChatMessageServer) ID() int32 {
	return 0x05
}

func (m *ChatMessageServer) Decode(r *encoding.Reader) error {
	r.String(&m.Message)
	r.Int64(&m.Timestamp)
	r.Int64(&m.Salt)

	var is bool
	r.Bool(&is)

	if is {
		m.Signature = make([]byte, 256)
		r.FixedByteArray(&m.Signature)
	}

	var count int32
	r.VarInt(&count)
	m.AcknowledgedMessages = make([]int64, count)

	for i := int32(0); i < count; i++ {
		r.Int64(&m.AcknowledgedMessages[i])
	}
	return nil
}

func (m ChatMessageServer) Encode(w *encoding.Writer) error {
	return w.String(m.Message)
}
