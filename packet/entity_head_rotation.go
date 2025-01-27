package packet

import "github.com/aimjel/minenet/protocol/encoding"

type EntityHeadRotation struct {
	EntityID int32
	HeadYaw  uint8
}

func (l EntityHeadRotation) ID() int32 {
	return 0x42
}

func (l *EntityHeadRotation) Decode(r *encoding.Reader) error {
	_ = r.VarInt(&l.EntityID)
	return r.Uint8(&l.HeadYaw)
}

func (l EntityHeadRotation) Encode(w *encoding.Writer) error {
	_ = w.VarInt(l.EntityID)
	return w.Uint8(l.HeadYaw)
}
