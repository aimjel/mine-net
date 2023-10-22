package packet

type BlockUpdate struct {
	Location int64
	BlockID  int32
}

func (c BlockUpdate) ID() int32 {
	return 0x0A
}

func (c *BlockUpdate) Decode(r *Reader) error {
	r.Int64(&c.Location)
	return r.VarInt(&c.BlockID)
}

func (c BlockUpdate) Encode(w Writer) error {
	w.Int64(c.Location)
	return w.VarInt(c.BlockID)
}
