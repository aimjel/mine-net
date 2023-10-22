package packet

type SetHeldItem struct {
	  Slot int8
}

func (m SetHeldItem) ID() int32 {
	return 0x4D
}

func (m *SetHeldItem) Decode(r *Reader) error {
	return r.Int8(&m.Slot)
}

func (m SetHeldItem) Encode(w Writer) error {
  return w.Int8(m.Slot)
}
