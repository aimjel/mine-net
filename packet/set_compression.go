package packet

type SetCompression struct {
	Threshold int32
}

func (s SetCompression) ID() int32 {
	return 0x03
}

func (s *SetCompression) Decode(r *Reader) error {
	return r.VarInt(&s.Threshold)
}

func (s SetCompression) Encode(w Writer) error {
	return w.VarInt(s.Threshold)
}
