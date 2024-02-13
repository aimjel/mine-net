package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type GameEvent struct {
	Event uint8
	Value float32
}

func (m GameEvent) ID() int32 {
	return 0x1F
}

func (m *GameEvent) Decode(r *encoding.Reader) error {
	r.Uint8(&m.Event)
	return r.Float32(&m.Value)
}

func (m GameEvent) Encode(w *encoding.Writer) error {
	w.Uint8(m.Event)
	return w.Float32(m.Value)
}
