package packet

type SpawnPlayer struct {
	EntityID   int32
	PlayerUUID [16]byte
	X, Y, Z    float64
	Yaw, Pitch byte
}

func (p SpawnPlayer) ID() int32 {
	return 0x03
}

func (p *SpawnPlayer) Decode(r *Reader) error {
	_ = r.VarInt(&p.EntityID)
	_ = r.UUID(&p.PlayerUUID)
	_ = r.Float64(&p.X)
	_ = r.Float64(&p.Y)
	_ = r.Float64(&p.Z)
	_ = r.Uint8(&p.Yaw)
	return r.Uint8(&p.Pitch)
}

func (p SpawnPlayer) Encode(w *Writer) error {
	_ = w.VarInt(p.EntityID)
	_ = w.UUID(p.PlayerUUID)
	_ = w.Float64(p.X)
	_ = w.Float64(p.Y)
	_ = w.Float64(p.Z)
	_ = w.Uint8(p.Yaw)
	return w.Uint8(p.Pitch)
}
