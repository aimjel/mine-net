package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type SetEquipment struct {
	EntityID int32
	Slot     int8
	Item     Slot
}

func (m SetEquipment) ID() int32 {
	return 0x55
}

func (m *SetEquipment) Decode(r *encoding.Reader) error {
	r.VarInt(&m.EntityID)
	r.Int8(&m.Slot)
	r.Bool(&m.Item.Present)
	if m.Item.Present {
		r.VarInt(&m.Item.Id)
		r.Int8(&m.Item.Count)
	}
	return nil
}

func (m SetEquipment) Encode(w *encoding.Writer) error {
	w.VarInt(m.EntityID)
	w.Int8(m.Slot)
	w.Bool(m.Item.Present)
	if m.Item.Present {
		w.VarInt(m.Item.Id)
		w.Int8(m.Item.Count)
		w.Nbt2(m.Item.Tag)
	}
	return nil
}
