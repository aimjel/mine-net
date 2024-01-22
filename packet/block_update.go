package packet

import (
	"github.com/aimjel/minecraft/protocol/encoding"
	"github.com/aimjel/minecraft/protocol/types"
)

type BlockUpdate struct {
	Location types.Position
	BlockID  int32
}

func (c BlockUpdate) ID() int32 {
	return 0x0A
}

func (c *BlockUpdate) Decode(r *encoding.Reader) error {
	r.Int64((*int64)(&c.Location))
	return r.VarInt(&c.BlockID)
}

func (c BlockUpdate) Encode(w *encoding.Writer) error {
	w.Int64(int64(c.Location))
	return w.VarInt(c.BlockID)
}
