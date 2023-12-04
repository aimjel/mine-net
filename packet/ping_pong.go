package packet

type Ping struct {
	Payload int64
}

func (p Ping) ID() int32 {
	return 0x01
}

func (p *Ping) Decode(r *Reader) error {
	return r.Int64(&p.Payload)
}

func (p Ping) Encode(w *Writer) error {
	return w.Int64(p.Payload)
}

type Pong struct {
	Payload int64
}

func (p Pong) ID() int32 {
	return 0x01
}

func (p *Pong) Decode(r *Reader) error {
	return r.Int64(&p.Payload)
}

func (p Pong) Encode(w *Writer) error {
	return w.Int64(p.Payload)
}
