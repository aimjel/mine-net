package packet

type HeldItemChange struct {
    Slot int16
}

func (c *HeldItemChange) ID() int32 {
    return 0x25
}

func (c *HeldItemChange) Encode(w Writer) error {
    return w.Int16(c.Slot)
}

func (c *HeldItemChange) Decode(r *Reader) error {
    return r.Int16(&c.Slot)
}

