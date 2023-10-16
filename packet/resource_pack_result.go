package packet

type ResourcePackResult struct {
	Result int32
}

func (m ResourcePackResult) ID() int32 {
	return 0x24
}

func (m *ResourcePackResult) Decode(r *Reader) error {
	return r.VarInt(&m.Result)
}

func (m ResourcePackResult) Encode(w Writer) error {
  return w.VarInt(m.Result)
}
