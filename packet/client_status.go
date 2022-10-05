package packet

type ClientStatus struct {
	ActionID int32
}

func (s ClientStatus) ID() int32 {
	return 0x04
}

func (s *ClientStatus) Encode(w Writer) error {
	return w.VarInt(s.ActionID)
}

func (s ClientStatus) Decode(r *Reader) error {
	return r.VarInt(&s.ActionID)
}
