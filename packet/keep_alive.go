package packet

type KeepAliveServer struct {
	PayloadID int64
}

func (a *KeepAliveServer) ID() int32 {
	return 0x12
}

func (a *KeepAliveServer) Encode(w *Writer) error {
	return w.Int64(a.PayloadID)
}

func (a *KeepAliveServer) Decode(r *Reader) error {
	return r.Int64(&a.PayloadID)
}

type KeepAliveClient struct {
	PayloadID int64
}

func (a *KeepAliveClient) ID() int32 {
	return 0x23
}

func (a *KeepAliveClient) Encode(w *Writer) error {
	return w.Int64(a.PayloadID)
}

func (a *KeepAliveClient) Decode(r *Reader) error {
	return r.Int64(&a.PayloadID)
}
