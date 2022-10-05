package packet

type EntityHeadLook struct {
	EntityID int32
	HeadYaw  uint8
}

func (l EntityHeadLook) ID() int32 {
	return 0x3E
}

func (l *EntityHeadLook) Decode(r *Reader) error {
	_ = r.VarInt(&l.EntityID)
	return r.Uint8(&l.HeadYaw)
}

func (l EntityHeadLook) Encode(w Writer) error {
	_ = w.VarInt(l.EntityID)
	return w.Uint8(l.HeadYaw)
}
