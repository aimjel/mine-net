package packet

type SetEntityVelocity struct {
	EntityID int32
	X, Y, Z int16
}

func (p SetEntityVelocity) ID() int32 {
	return 0x54
}

func (p *SetEntityVelocity) Decode(r *Reader) error {
	_ = r.VarInt(&p.EntityID)
	_ = r.Int16(&p.X)
	_ = r.Int16(&p.Y)
	return r.Int16(&p.Z)
}
func (p SetEntityVelocity) Encode(w Writer) error {
	_ = w.VarInt(p.EntityID)
	_ = w.Int16(p.X)
	_ = w.Int16(p.Y)
	return w.Int16(p.Z)
}
