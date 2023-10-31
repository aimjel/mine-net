package packet

type DeleteMessage struct {
  MessageID int32
  Signature []byte
}

func (m DeleteMessage) ID() int32 {
	return 0x19
}

func (m *DeleteMessage) Decode(r *Reader) error {
	r.VarInt(&m.MessageID)
  m.MessageID--
  if m.MessageID == -1 {
    m.Signature = make([]byte, 256)
    r.FixedByteArray(&m.Signature)
  }
  return nil
}

func (m DeleteMessage) Encode(w Writer) error {
  w.VarInt(m.MessageID+1)
  if m.MessageID + 1 == 0 {
    w.FixedByteArray(m.Signature)
  }
  return nil
}
