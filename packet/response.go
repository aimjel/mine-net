package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type Response struct {
	JSON []byte
}

func (r Response) ID() int32 {
	return 0x00
}

func (r *Response) Decode(rd *encoding.Reader) error {
	return rd.ByteArray(&r.JSON)
}

func (r Response) Encode(w *encoding.Writer) error {
	return w.ByteArray(r.JSON)
}
