package packet

type TeleportToEntityServer struct {
	Player [16]byte
}

func (s TeleportToEntityServer) ID() int32 {
	return 0x30
}

func (s *TeleportToEntityServer) Decode(r *Reader) error {
	return r.UUID(&s.Player)
}

func (s TeleportToEntityServer) Encode(w Writer) error {
	return w.UUID(s.Player)
}
