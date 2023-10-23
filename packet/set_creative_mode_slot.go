package packet

type SetCreativeModeSlot struct {
	Slot int16
  	ClickedItem Slot
  //TODO clicked item tag
}

type Slot struct {
	Count int8   
	Id    int32
	Tag SlotTag
}

type Enchantment struct {
	Id    string
	Level int16
}

type SlotTag struct {
	Damage       int32    
	RepairCost   int32       
	Enchantments []Enchantment
}

func (c SetCreativeModeSlot) ID() int32 {
	return 0x2B
}

func (c *SetCreativeModeSlot) Decode(r *Reader) error {
	r.Int16(&c.Slot)
	var present bool
	r.Bool(&present)
	if present {
		r.VarInt(&c.ClickedItem.Id)
		r.Int8(&c.ClickedItem.Count)
	}
	return nil
}

func (c SetCreativeModeSlot) Encode(w Writer) error {
	return w.Int16(c.Slot)
}
