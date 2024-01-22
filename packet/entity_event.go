package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type EntityEvent struct {
	EntityID int32
	Status   int8
}

func (c EntityEvent) ID() int32 {
	return 0x1C
}

func (c *EntityEvent) Decode(r *encoding.Reader) error {
	r.Int32(&c.EntityID)
	return r.Int8(&c.Status)
}

func (c EntityEvent) Encode(w *encoding.Writer) error {
	w.Int32(c.EntityID)
	return w.Int8(c.Status)
}
