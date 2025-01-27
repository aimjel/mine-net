package packet

import "github.com/aimjel/minenet/protocol/encoding"

type PaddleBoat struct {
	LeftPaddleTurning  bool
	RightPaddleTurning bool
}

func (p PaddleBoat) ID() int32 {
	return 0x19
}

func (p *PaddleBoat) Decode(r *encoding.Reader) error {
	r.Bool(&p.LeftPaddleTurning)
	return r.Bool(&p.RightPaddleTurning)
}

func (p *PaddleBoat) Encode(w *encoding.Writer) error {
	w.Bool(p.LeftPaddleTurning)
	return w.Bool(p.RightPaddleTurning)
}
