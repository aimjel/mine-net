package packet

type UseItem struct {
	Hand int32
}

func (i *UseItem) ID() int32 {
	return 0x2f
}

func (i *UseItem) Encode(w Writer) error {
	return nil
}

func (i *UseItem) Decode(r *Reader) error {
	return r.VarInt(&i.Hand)
}
