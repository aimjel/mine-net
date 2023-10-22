package packet

type PlayerAbilitiesServer struct {
	Flags byte
}

func (p PlayerAbilitiesServer) ID() int32 {
	return 0x1C
}

func (p *PlayerAbilitiesServer) Decode(r *Reader) error {
	return r.Uint8(&p.Flags)
}
func (p PlayerAbilitiesServer) Encode(w Writer) error {
	return w.Uint8(p.Flags)
}
