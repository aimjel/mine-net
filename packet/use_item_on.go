package packet

type UseItemOnServer struct {
	Hand int32
  Location int64
  Face int32
  CursorPositionX, CursorPositionY, CursorPositionZ float32
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
  return r.Float32(&m.CursorPositionZ)
}

func (m UseItemOnServer) Encode(w Writer) error {
	w.VarInt(m.Hand)
  w.Int64(m.Location)
  w.VarInt(m.Face)
  w.Float32(m.CursorPositionX)
  w.Float32(m.CursorPositionY)
  return w.Float32(m.CursorPositionZ)
}
