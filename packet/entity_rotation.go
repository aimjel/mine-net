package packet

import "github.com/aimjel/minenet/protocol/encoding"

type EntityRotation struct {
	EntityID   int32
	Yaw, Pitch uint8
	OnGround   bool
}

func (r EntityRotation) ID() int32 {
	return 0x2D
}

func (r *EntityRotation) Decode(rd *encoding.Reader) error {
	_ = rd.VarInt(&r.EntityID)
	_ = rd.Uint8(&r.Yaw)
	_ = rd.Uint8(&r.Pitch)
	return rd.Bool(&r.OnGround)
}

func (r EntityRotation) Encode(w *encoding.Writer) error {
	_ = w.VarInt(r.EntityID)
	_ = w.Uint8(r.Yaw)
	_ = w.Uint8(r.Pitch)
	return w.Bool(r.OnGround)
}
