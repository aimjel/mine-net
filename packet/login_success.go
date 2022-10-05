package packet

type LoginSuccess struct {
	UUID     [16]byte
	Username string
}

func (s LoginSuccess) ID() int32 {
	return 0x02
}

func (s *LoginSuccess) Decode(r *Reader) error {
	_ = r.UUID(&s.UUID)
	return r.String(&s.Username)
}

func (s LoginSuccess) Encode(w Writer) error {
	_ = w.UUID(s.UUID)
	return w.String(s.Username)
}
