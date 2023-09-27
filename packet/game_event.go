package packet

type GameEvent struct {
  Event uint8
  Value float32
}

func (m GameEvent) ID() int32 {
	return 0x20
}

func (m *GameEvent) Decode(r *Reader) error {
 	r.Uint8(m.Event)
	r.Float32(m.Value)
}

func (m GameEvent) Encode(w Writer) error {
	w.Uint8(m.Event)
	return w.Float32(m.Value)
}
