package packet

type InteractServer struct {
	EntityID int32
	Type     int32
	TargetX  float32
	TargetY  float32
	TargetZ  float32
	Hand     int32
	Sneaking bool
}

func (p InteractServer) ID() int32 {
	return 0x10
}

func (p *InteractServer) Decode(r *Reader) error {
	_ = r.VarInt(&p.EntityID)
	_ = r.VarInt(&p.Type)
	if p.Type == 2 {
		_ = r.Float32(&p.TargetX)
		_ = r.Float32(&p.TargetY)
		_ = r.Float32(&p.TargetZ)
		_ = r.VarInt(&p.Hand)
	}
	if p.Type == 0 {
		_ = r.VarInt(&p.Hand)
	}

	return r.Bool(&p.Sneaking)
}

func (p InteractServer) Encode(w *Writer) error {
	_ = w.VarInt(p.EntityID)
	_ = w.VarInt(p.Type)
	if p.Type == 2 {
		_ = w.Float32(p.TargetX)
		_ = w.Float32(p.TargetY)
		_ = w.Float32(p.TargetZ)
		_ = w.VarInt(p.Hand)
	}
	if p.Type == 0 {
		_ = w.VarInt(p.Hand)
	}

	return w.Bool(p.Sneaking)
}
