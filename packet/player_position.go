package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type PlayerPosition struct {
	X, FeetY, Z float64

	OnGround bool
}

func (p PlayerPosition) ID() int32 {
	return 0x14
}

func (p *PlayerPosition) Decode(r *encoding.Reader) error {
	_ = r.Float64(&p.X)
	_ = r.Float64(&p.FeetY)
	_ = r.Float64(&p.Z)

	return r.Bool(&p.OnGround)
}

func (p PlayerPosition) Encode(w *encoding.Writer) error {
	_ = w.Float64(p.X)
	_ = w.Float64(p.FeetY)
	_ = w.Float64(p.Z)

	return w.Bool(p.OnGround)
}
