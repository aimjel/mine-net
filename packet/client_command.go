package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type ClientCommandServer struct {
	ActionID int32
}

func (m ClientCommandServer) ID() int32 {
	return 0x07
}

func (m *ClientCommandServer) Decode(r *encoding.Reader) error {
	return r.VarInt(&m.ActionID)
}

func (m ClientCommandServer) Encode(w *encoding.Writer) error {
	return w.VarInt(m.ActionID)
}
