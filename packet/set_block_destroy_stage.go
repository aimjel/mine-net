package packet

type SetBlockDestroyStage struct {
	EntityID int32
	Location  uint64
  DestroyStage byte
}

func (c SetBlockDestroyStage) ID() int32 {
	return 0x07
}

func (c *SetBlockDestroyStage) Decode(r *Reader) error {
	r.VarInt(&c.EntityID)
  r.Uint64(&c.Location)
  return r.Uint8(&c.DestroyStage)
}

func (c SetBlockDestroyStage) Encode(w Writer) error {
	w.VarInt(c.EntityID)
  w.Uint64(c.Location)
  return w.Uint8(c.DestroyStage)
}
