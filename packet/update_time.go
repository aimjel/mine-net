package packet

type UpdateTime struct {
	WorldAge  int64
	TimeOfDay int64
}

func (c UpdateTime) ID() int32 {
	return 0x5E
}

func (c *UpdateTime) Decode(r *Reader) error {
	r.Int64(&c.WorldAge)
	return r.Int64(&c.TimeOfDay)
}

func (c UpdateTime) Encode(w *Writer) error {
	w.Int64(c.WorldAge)
	return w.Int64(c.TimeOfDay)
}
