package packet

type PlayerBlockPlacement struct {
	Hand            int32
	Location        uint64
	Face            int32
	CursorPositionX float32
	CursorPositionY float32
	CursorPositionZ float32
	InsideBlock     bool
}

func (p *PlayerBlockPlacement) ID() int32 {
	return 0x2e
}

func (p *PlayerBlockPlacement) Decode(r *Reader) error {
	_ = r.VarInt(&p.Hand)
	_ = r.Uint64(&p.Location)
	_ = r.VarInt(&p.Face)
	_ = r.Float32(&p.CursorPositionX)
	_ = r.Float32(&p.CursorPositionY)
	_ = r.Float32(&p.CursorPositionZ)
	return r.Bool(&p.InsideBlock)
}

func (p *PlayerBlockPlacement) Encode(w Writer) error {
	return nil
}
