package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type SwingArmServer struct {
	Hand int32
}

func (s SwingArmServer) ID() int32 {
	return 0x2F
}

func (s *SwingArmServer) Decode(r *encoding.Reader) error {
	return r.VarInt(&s.Hand)
}

func (s SwingArmServer) Encode(w *encoding.Writer) error {
	return w.VarInt(s.Hand)
}
