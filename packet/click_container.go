package packet

type ClickContainer struct {
	WindowID    uint8
	StateID     int32
	Slot        int16
	Button      int8
	Mode        int32
	Slots       map[int16]Slot
	CarriedItem Slot
}

func (c ClickContainer) ID() int32 {
	return 0x0B
}

func (c *ClickContainer) Decode(r *Reader) error {
	c.Slots = make(map[int16]Slot)
	r.Uint8(&c.WindowID)
	r.VarInt(&c.StateID)
	r.Int16(&c.Slot)
	r.Int8(&c.Button)
	r.VarInt(&c.Mode)

	var size int32
	r.VarInt(&size)

	for i := int32(0); i < size; i++ {
		var slotNum int16
		r.Int16(&slotNum)

		slot := c.Slots[slotNum]

		r.Bool(&slot.Present)

		if slot.Present {
			r.VarInt(&slot.Id)
			r.Int8(&slot.Count)

			c.Slots[slotNum] = slot
		}
	}

	return nil
}

func (c ClickContainer) Encode(w Writer) error {
	return w.Uint8(c.WindowID)
}
