package packet

import "github.com/aimjel/minenet/protocol/encoding"

type PlayerAbilities struct {
	Flags               byte
	FlyingSpeed         float32
	FieldOfViewModifier float32
}

func (p PlayerAbilities) ID() int32 {
	return 0x34
}

func (p *PlayerAbilities) Decode(r *encoding.Reader) error {
	_ = r.Uint8(&p.Flags)
	_ = r.Float32(&p.FlyingSpeed)
	return r.Float32(&p.FieldOfViewModifier)
}
func (p PlayerAbilities) Encode(w *encoding.Writer) error {
	_ = w.Uint8(p.Flags)
	_ = w.Float32(p.FlyingSpeed)
	return w.Float32(p.FieldOfViewModifier)
}
