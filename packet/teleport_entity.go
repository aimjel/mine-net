package packet

type TeleportEntity struct {
	EntityID   int32
	X, Y, Z    float64
	Yaw, Pitch byte
	OnGround   bool
}

func (r TeleportEntity) ID() int32 {
	return 0x68
}

func (r *TeleportEntity) Decode(rd *Reader) error {
	_ = rd.VarInt(&r.EntityID)
	_ = rd.Float64(&r.X)
	_ = rd.Float64(&r.Y)
	_ = rd.Float64(&r.Z)
	_ = rd.Uint8(&r.Yaw)
	_ = rd.Uint8(&r.Pitch)
	return rd.Bool(&r.OnGround)
}

func (r TeleportEntity) Encode(w *Writer) error {
	_ = w.VarInt(r.EntityID)
	_ = w.Float64(r.X)
	_ = w.Float64(r.Y)
	_ = w.Float64(r.Z)
	_ = w.Uint8(r.Yaw)
	_ = w.Uint8(r.Pitch)
	return w.Bool(r.OnGround)
}
