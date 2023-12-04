package packet

type LoginStart struct {
	Name string
	UUID [16]byte
}

func (s LoginStart) ID() int32 {
	return 0x00
}

func (s *LoginStart) Decode(r *Reader) error {
	_ = r.String(&s.Name)
	var hasUUID bool
	_ = r.Bool(&hasUUID)
	if hasUUID {
		return r.UUID(&s.UUID)
	}

	return nil
}

func (s LoginStart) Encode(w *Writer) error {
	_ = w.String(s.Name)
	var hasUUID bool
	_ = w.Bool(hasUUID)
	if hasUUID {
		return w.UUID(s.UUID)
	}

	return nil
}
