package packet

type LoginStart struct {
	Name string
}

func (s LoginStart) ID() int32 {
	return 0x00
}

func (s *LoginStart) Decode(r *Reader) error {
	return r.String(&s.Name)
}

func (s LoginStart) Encode(w Writer) error {
	return w.String(s.Name)
}
