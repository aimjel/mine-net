package packet

type PlayerPositionLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      int8
	TeleportID int32
}

func (l *PlayerPositionLook) ID() int32 {
	return 0x3c
}

func (l *PlayerPositionLook) Decode(r *Reader) error {
	_ = r.Float64(&l.X)
	_ = r.Float64(&l.Y)
	_ = r.Float64(&l.Z)
	_ = r.Float32(&l.Yaw)
	_ = r.Float32(&l.Pitch)
	_ = r.Int8(&l.Flags)
	return r.VarInt(&l.TeleportID)
}

func (l PlayerPositionLook) Encode(w *Writer) error {
	_ = w.Float64(l.X)
	_ = w.Float64(l.Y)
	_ = w.Float64(l.Z)
	_ = w.Float32(l.Yaw)
	_ = w.Float32(l.Pitch)
	_ = w.Int8(l.Flags)
	return w.VarInt(l.TeleportID)
}
