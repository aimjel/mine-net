package packet

type SetCreativeModeSlot struct {
	Slot int16
  //ClickedItem Slot
  //TODO clicked item
}

func (c SetCreativeModeSlot) ID() int32 {
	return 0x2B
}

func (c *SetCreativeModeSlot) Decode(r *Reader) error {
	return r.Int16(&c.Slot)
}

func (c SetCreativeModeSlot) Encode(w Writer) error {
	return w.Int16(c.Slot)
}
