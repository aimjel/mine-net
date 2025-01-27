package packet

import "github.com/aimjel/minenet/protocol/encoding"

type EntityPositionRotation struct {
	EntityID   int32
	X, Y, Z    int16
	Yaw, Pitch byte
	OnGround   bool
}

func (r EntityPositionRotation) ID() int32 {
	return 0x2C
}

func (r *EntityPositionRotation) Decode(rd *encoding.Reader) error {
	_ = rd.VarInt(&r.EntityID)
	_ = rd.Int16(&r.X)
	_ = rd.Int16(&r.Y)
	_ = rd.Int16(&r.Z)
	_ = rd.Uint8(&r.Yaw)
	_ = rd.Uint8(&r.Pitch)
	return rd.Bool(&r.OnGround)
}

func (r EntityPositionRotation) Encode(w *encoding.Writer) error {
	_ = w.VarInt(r.EntityID)
	_ = w.Int16(r.X)
	_ = w.Int16(r.Y)
	_ = w.Int16(r.Z)
	_ = w.Uint8(r.Yaw)
	_ = w.Uint8(r.Pitch)
	return w.Bool(r.OnGround)
}
