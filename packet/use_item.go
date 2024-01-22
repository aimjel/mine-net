package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type UseItem struct {
	Hand int32
}

func (i *UseItem) ID() int32 {
	return 0x2f
}

func (i *UseItem) Encode(w *encoding.Writer) error {
	return nil
}

func (i *UseItem) Decode(r *encoding.Reader) error {
	return r.VarInt(&i.Hand)
}
