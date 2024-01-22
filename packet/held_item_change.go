package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type HeldItemChange struct {
	Slot int16
}

func (c *HeldItemChange) ID() int32 {
	return 0x25
}

func (c *HeldItemChange) Encode(w *encoding.Writer) error {
	return w.Int16(c.Slot)
}

func (c *HeldItemChange) Decode(r *encoding.Reader) error {
	return r.Int16(&c.Slot)
}
