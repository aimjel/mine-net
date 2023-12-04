package packet

type DestroyEntities struct {
	EntityIds []int32
}

func (d DestroyEntities) ID() int32 {
	return 0x3E
}

func (d *DestroyEntities) Decode(r *Reader) error {
	panic("implement me")
}

func (d DestroyEntities) Encode(w *Writer) error {
	return w.VarIntArray(d.EntityIds)
}
