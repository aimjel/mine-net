package packet

import "github.com/aimjel/minenet/protocol/encoding"

type HurtAnimation struct {
	EntityID int32
	Yaw      float32
}

func (l HurtAnimation) ID() int32 {
	return 0x21
}

func (l *HurtAnimation) Decode(r *encoding.Reader) error {
	r.VarInt(&l.EntityID)
	return r.Float32(&l.Yaw)
}

func (l HurtAnimation) Encode(w *encoding.Writer) error {
	w.VarInt(l.EntityID)
	return w.Float32(l.Yaw)
}
