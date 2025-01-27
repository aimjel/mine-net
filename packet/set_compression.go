package packet

import "github.com/aimjel/minenet/protocol/encoding"

type SetCompression struct {
	Threshold int32
}

func (s SetCompression) ID() int32 {
	return 0x03
}

func (s *SetCompression) Decode(r *encoding.Reader) error {
	return r.VarInt(&s.Threshold)
}

func (s SetCompression) Encode(w *encoding.Writer) error {
	return w.VarInt(s.Threshold)
}
