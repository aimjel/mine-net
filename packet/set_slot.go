package packet

type SetContainerSlot struct {
	WindowID  int8
	StateID   int32
	Slot      int16
	Data Slot
}

func (s *SetContainerSlot) ID() int32 {
	return 0x14
}

func (s *SetContainerSlot) Encode(w Writer) error {
	_ = w.Int8(s.WindowID)
	_ = w.VarInt(s.StateID)
	_ = w.Int16(s.Slot)
	_ = w.Bool(s.Data.Present)
	if s.Present {
		_ = w.VarInt(s.Data.Id)
		_ = w.Int8(s.Data.Count)
		return w.Nbt2(s.Data.Tag)
	}

	return nil
}
