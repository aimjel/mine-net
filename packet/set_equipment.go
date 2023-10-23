package packet

type SetEquipment struct {
  	 EntityID int32
     Slot int8
     Item Slot
}

func (m SetEquipment) ID() int32 {
	return 0x55
}

func (m *SetEquipment) Decode(r *Reader) error {
	r.VarInt(&m.EntityID)
  r.Int8(&m.Slot)
  r.Bool(&c.Item.Present)
	if c.Item.Present {
		r.VarInt(&c.Item.Id)
		r.Int8(&c.Item.Count)
	}
  return nil
}

func (m SetEquipment) Encode(w Writer) error {
 	w.VarInt(m.EntityID)
  w.Int8(m.Slot)
  w.Bool(c.Item.Present)
	if c.Item.Present {
		w.VarInt(c.Item.Id)
		w.Int8(c.Item.Count)
    w.Nbt2(c.Item.Tag)
	}
  return nil
}
