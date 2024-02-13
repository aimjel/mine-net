package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type PlayerAbilitiesServer struct {
	Flags byte
}

func (p PlayerAbilitiesServer) ID() int32 {
	return 0x1C
}

func (p *PlayerAbilitiesServer) Decode(r *encoding.Reader) error {
	return r.Uint8(&p.Flags)
}
func (p PlayerAbilitiesServer) Encode(w *encoding.Writer) error {
	return w.Uint8(p.Flags)
}
