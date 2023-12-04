package packet

type SetContainerContent struct {
	WindowID uint8
	StateID  int32
	Slots    []Slot
}

func (m SetContainerContent) ID() int32 {
	return 0x12
}

func (m *SetContainerContent) Decode(r *Reader) error {
	//todo reader
	return nil
}

func (m SetContainerContent) Encode(w *Writer) error {
	w.Uint8(m.WindowID)
	w.VarInt(m.StateID)
	w.VarInt(int32(len(m.Slots)))
	for _, s := range m.Slots {
		if !s.Present {
			w.Bool(false)
			continue
		}
		w.Bool(true)
		w.VarInt(s.Id)
		w.Int8(s.Count)
		w.Nbt2(s.Tag)
	}
	w.Bool(false)
	return nil
}
