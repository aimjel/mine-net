package packet

import (
	"github.com/aimjel/minenet/protocol/encoding"
	"github.com/aimjel/minenet/protocol/types"
)

type PlayerActionServer struct {
	Status   int32
	Location types.Position
	Face     int8
	Sequence int32
}

func (m PlayerActionServer) ID() int32 {
	return 0x1D
}

func (m *PlayerActionServer) Decode(r *encoding.Reader) error {
	r.VarInt(&m.Status)
	r.Int64((*int64)(&m.Location))
	r.Int8(&m.Face)
	return r.VarInt(&m.Sequence)
}

func (m PlayerActionServer) Encode(w *encoding.Writer) error {
	w.VarInt(m.Status)
	w.Int64(int64(m.Location))
	w.Int8(m.Face)
	return w.VarInt(m.Sequence)
}
