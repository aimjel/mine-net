package packet

type EntityPosition struct {
	EntityID int32
	X, Y, Z  int16
	OnGround bool
}

func (p EntityPosition) ID() int32 {
	return 0x2B
}

func (p *EntityPosition) Decode(r *Reader) error {
	_ = r.VarInt(&p.EntityID)
	_ = r.Int16(&p.X)
	_ = r.Int16(&p.Y)
	_ = r.Int16(&p.Z)
	return r.Bool(&p.OnGround)
}
func (p EntityPosition) Encode(w Writer) error {
	_ = w.VarInt(p.EntityID)
	_ = w.Int16(p.X)
	_ = w.Int16(p.Y)
	_ = w.Int16(p.Z)
	return w.Bool(p.OnGround)
}
