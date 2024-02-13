package packet

import (
	"github.com/aimjel/minecraft/protocol/encoding"
	"github.com/aimjel/minecraft/protocol/types"
)

type SetBlockDestroyStage struct {
	EntityID     int32
	Location     types.Position
	DestroyStage byte
}

func (c SetBlockDestroyStage) ID() int32 {
	return 0x07
}

func (c *SetBlockDestroyStage) Decode(r *encoding.Reader) error {
	r.VarInt(&c.EntityID)
	r.Int64((*int64)(&c.Location))
	return r.Uint8(&c.DestroyStage)
}

func (c SetBlockDestroyStage) Encode(w *encoding.Writer) error {
	w.VarInt(c.EntityID)
	w.Int64(int64(c.Location))
	return w.Uint8(c.DestroyStage)
}
