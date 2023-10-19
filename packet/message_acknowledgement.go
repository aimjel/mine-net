package packet

type MessageAcknowledgment struct {
	MessageCount int32
}

func (m MessageAcknowledgment) ID() int32 {
	return 0x03
}

func (m *MessageAcknowledgment) Decode(r *Reader) error {
  return r.VarInt(&m.MessageCount)
}

func (m MessageAcknowledgment) Encode(w Writer) error {
	return w.VarInt(m.MessageCount)
}
