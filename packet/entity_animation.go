package packet

type EntityAnimation struct {
	EntityID  int32
	Animation uint8
}

func (a EntityAnimation) ID() int32 {
	return 0x04
}

func (a *EntityAnimation) Decode(r *Reader) error {
	r.VarInt(&a.EntityID)
	return r.Uint8(&a.Animation)
}

func (a EntityAnimation) Encode(w Writer) error {
	w.VarInt(a.EntityID)
	return w.Uint8(a.Animation)
}
