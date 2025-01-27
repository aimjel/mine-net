package packet

import "github.com/aimjel/minenet/protocol/encoding"

type DestroyEntities struct {
	EntityIds []int32
}

func (d DestroyEntities) ID() int32 {
	return 0x3E
}

func (d *DestroyEntities) Decode(r *encoding.Reader) error {
	panic("implement me")
}

func (d DestroyEntities) Encode(w *encoding.Writer) error {
	return w.VarIntArray(d.EntityIds)
}
