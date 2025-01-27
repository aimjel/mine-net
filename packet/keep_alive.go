package packet

import "github.com/aimjel/minenet/protocol/encoding"

type KeepAliveServer struct {
	PayloadID int64
}

func (a *KeepAliveServer) ID() int32 {
	return 0x12
}

func (a *KeepAliveServer) Encode(w *encoding.Writer) error {
	return w.Int64(a.PayloadID)
}

func (a *KeepAliveServer) Decode(r *encoding.Reader) error {
	return r.Int64(&a.PayloadID)
}

type KeepAliveClient struct {
	PayloadID int64
}

func (a *KeepAliveClient) ID() int32 {
	return 0x23
}

func (a *KeepAliveClient) Encode(w *encoding.Writer) error {
	return w.Int64(a.PayloadID)
}

func (a *KeepAliveClient) Decode(r *encoding.Reader) error {
	return r.Int64(&a.PayloadID)
}
