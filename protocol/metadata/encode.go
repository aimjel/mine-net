package metadata

import (
	"fmt"
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

func encode(w *encoding.Writer, t any) error {
	switch v := t.(type) {

	case byte:
		//type id, value
		fmt.Println("byte type", 0, v)
		return w.FixedByteArray([]byte{0, v})

	case int32:
		_ = w.Uint8(1)
		return w.VarInt(v)

	case *chat.Message:
		fmt.Println("encoding chat message")
		_ = w.Uint8(6)
		_ = w.Bool(v != nil)
		if v != nil {
			return w.String(v.String())
		}
		return nil

	case bool:
		_ = w.Uint8(8)
		return w.Bool(v)

	case pose:
		_ = w.Uint8(20)
		fmt.Println("pose type", 20, v)
		return w.VarInt(int32(v))

	default:
		panic("unknown metadata type")
	}
}
