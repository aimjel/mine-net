package packet

type SetCreativeModeSlot struct {
	Slot int16
  	ClickedItem Slot
}

type Slot struct {
	Present bool
	Count int8   
	Id    int32
	Tag SlotTag
}

type Enchantment struct {
	Id    string `nbt:"id"`
	Level int16  `nbt:"lvl"`
}

type SlotTag struct {
	Damage       int32         `nbt:"Damage"`
	RepairCost   int32         `nbt:"RepairCost"`
	Enchantments []Enchantment `nbt:"Enchantments"`
}

func (c SetCreativeModeSlot) ID() int32 {
	return 0x2B
}

func (c *SetCreativeModeSlot) Decode(r *Reader) error {
	r.Int16(&c.Slot)
	r.Bool(&c.ClickedItem.Present)
	if c.ClickedItem.Present {
		r.VarInt(&c.ClickedItem.Id)
		r.Int8(&c.ClickedItem.Count)
		r.Nbt(&c.ClickedItem.Tag)
	}
	return nil
}

func (c SetCreativeModeSlot) Encode(w Writer) error {
	return w.Int16(c.Slot)
}
