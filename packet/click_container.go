package packet

type ClickContainer struct {
	WindowID uint8
  StateID int32
  Slot int16
  Button int8
  Mode int32
  Slots map[int16]Slot
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
  r.VarInt(&c.Size)

  for i := int32(0); i < size; i++ {
    var (
       slot int16
    )
    r.Int16(&slot)
    r.Bool(&c.Slots[slot].Present)
    if c.Slots[slot].Present {
      r.VarInt(&c.Slots[slot].Id)
		  r.Int8(&c.Slots[slot].Count)

      // pause
    }
  }

	return nil
}

func (c ClickContainer) Encode(w Writer) error {
	return w.Uint8(c.WindowID)
}
