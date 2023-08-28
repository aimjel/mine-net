package packet

type KeepAlive struct {
	PayloadID int64
}

func (a *KeepAlive) ID() int32 {
	return 0x23
}

func (a *KeepAlive) Encode(w Writer) error {
	return w.Int64(a.PayloadID)
}

func (a *KeepAlive) Decode(r *Reader) error {
	return r.Int64(&a.PayloadID)
}
