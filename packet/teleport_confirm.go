package packet

type TeleportConfirm struct {
	TeleportID int32
}

func (t TeleportConfirm) ID() int32 {
	return 0x00
}

func (t *TeleportConfirm) Decode(r *Reader) error {
	return r.VarInt(&t.TeleportID)
}

func (t TeleportConfirm) Encode(w *Writer) error {
	return w.VarInt(t.TeleportID)
}
