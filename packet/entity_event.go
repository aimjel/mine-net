package packet

type EntityEvent struct {
   EntityID int32
   Status byte
}

func (c EntityEvent) ID() int32 {
	return 0x1C
}

func (c *EntityEvent) Decode(r *Reader) error {
  r.VarInt(&c.EntityID)
	return r.Uint8(&c.Status)
}

func (c EntityEvent) Encode(w Writer) error {
  w.VarInt(c.EntityID)
	return w.Uint8(c.Status)
}
