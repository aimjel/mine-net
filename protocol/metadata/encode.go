package metadata

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

func encode(w *encoding.Writer, t any) error {
	switch v := t.(type) {

	case byte:
		//type id, value
		return w.FixedByteArray([]byte{0, v})

	case int32:
		_ = w.Uint8(1)
		return w.VarInt(v)

	case *chat.Message:
		_ = w.Uint8(6)
		_ = w.Bool(v != nil)
		if v != nil {
			return w.String(v.String())
		}
		return nil

	case bool:
		_ = w.Uint8(8)
		return w.Bool(v)

	case Pose:
		_ = w.Uint8(20)
		return w.VarInt(v)

	default:
		panic("unknown metadata type")
	}
}
