package packet

import "github.com/aimjel/minenet/protocol/encoding"

type PlayerCommandServer struct {
	EntityID  int32
	ActionID  int32
	JumpBoost int32
}

func (m PlayerCommandServer) ID() int32 {
	return 0x1E
}

func (m *PlayerCommandServer) Decode(r *encoding.Reader) error {
	r.VarInt(&m.EntityID)
	r.VarInt(&m.ActionID)
	return r.VarInt(&m.JumpBoost)
}

func (m PlayerCommandServer) Encode(w *encoding.Writer) error {
	w.VarInt(m.EntityID)
	w.VarInt(m.ActionID)
	return w.VarInt(m.JumpBoost)
}
