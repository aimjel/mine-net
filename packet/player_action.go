package packet

type PlayerActionServer struct {
  Status int32
  Location uint64
  Face int8
  Sequence int32
}

func (m PlayerActionServer) ID() int32 {
	return 0x1D
}

func (m *PlayerActionServer) Decode(r *Reader) error {
	r.VarInt(&m.Status)
	r.Uint64(&m.Location)
  r.Int8(&m.Face)
	return r.VarInt(&m.Sequence)
}

func (m PlayerActionServer) Encode(w Writer) error {
	w.VarInt(m.Status)
	w.Uint64(m.Location)
  w.Int8(m.Face)
	return w.VarInt(m.Sequence)
}
