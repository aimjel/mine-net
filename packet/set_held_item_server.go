package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type SetHeldItemServer struct {
	Slot int16
}

func (m SetHeldItemServer) ID() int32 {
	return 0x28
}

func (m *SetHeldItemServer) Decode(r *encoding.Reader) error {
	return r.Int16(&m.Slot)
}

func (m SetHeldItemServer) Encode(w *encoding.Writer) error {
	return w.Int16(m.Slot)
}
