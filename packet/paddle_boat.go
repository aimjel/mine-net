package packet

type PaddleBoat struct {
	LeftPaddleTurning  bool
	RightPaddleTurning bool
}

func (p PaddleBoat) ID() int32 {
	return 0x19
}

func (p *PaddleBoat) Decode(r *Reader) error {
	r.Bool(&p.LeftPaddleTurning)
	return r.Bool(&p.RightPaddleTurning)
}

func (p *PaddleBoat) Encode(w *Writer) error {
	w.Bool(p.LeftPaddleTurning)
	return w.Bool(p.RightPaddleTurning)
}
