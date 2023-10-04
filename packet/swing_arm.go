package packet

type SwingArmServer struct {
	Hand int32
}

func (s SwingArmServer) ID() int32 {
	return 0x2F
}

func (s *SwingArmServer) Decode(r *Reader) error {
  return r.VarInt(&s.Hand)
}

func (s SwingArmServer) Encode(w Writer) error {
	return w.VarInt(s.Hand)
}
