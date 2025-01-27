package packet

import "github.com/aimjel/minenet/protocol/encoding"

type TeleportConfirm struct {
	TeleportID int32
}

func (t TeleportConfirm) ID() int32 {
	return 0x00
}

func (t *TeleportConfirm) Decode(r *encoding.Reader) error {
	return r.VarInt(&t.TeleportID)
}

func (t TeleportConfirm) Encode(w *encoding.Writer) error {
	return w.VarInt(t.TeleportID)
}
