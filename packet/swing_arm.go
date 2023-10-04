package packet

type SwingArm struct {
	Hand int32
}

func (s SwingArm) ID() int32 {
	return 0x2F
}

func (s *SwingArm) Decode(r *Reader) error {
  return r.VarInt(&s.Hand)
}

func (s SwingArm) Encode(w Writer) error {
	return w.VarInt(s.Hand)
}
