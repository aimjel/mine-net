package packet

type UseItemOnServer struct {
	Hand int32
  Location int64
  Face int32
  CursorPositionX, CursorPositionY, CursorPositionZ float32
	InsideBlock bool
	Sequence int32
}

func (m UseItemOnServer) ID() int32 {
	return 0x31
}

func (m *UseItemOnServer) Decode(r *Reader) error {
	  r.VarInt(&m.Hand)
	  r.Int64(&m.Location)
	  r.VarInt(&m.Face)
	  r.Float32(&m.CursorPositionX)
	  r.Float32(&m.CursorPositionY)
	  r.Float32(&m.CursorPositionZ)
	  r.Bool(&m.InsideBlock)
	  return r.VarInt(&m.Sequence)
}

func (m UseItemOnServer) Encode(w Writer) error {
	w.VarInt(m.Hand)
  w.Int64(m.Location)
  w.VarInt(m.Face)
  w.Float32(m.CursorPositionX)
  w.Float32(m.CursorPositionY)
  w.Float32(m.CursorPositionZ)
		  w.Bool(m.InsideBlock)
	  return w.VarInt(m.Sequence)
}
