package packet

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

type DisguisedChatMessage struct {
	Message      chat.Message
	ChatType     int32
	ChatTypeName chat.Message
	TargetName   *chat.Message
}

func (m DisguisedChatMessage) ID() int32 {
	return 0x1B
}

func (m *DisguisedChatMessage) Decode(r *encoding.Reader) error {
	return NotImplemented
}

func (m DisguisedChatMessage) Encode(w *encoding.Writer) error {
	w.String(m.Message.String())
	w.VarInt(m.ChatType)
	w.String(m.ChatTypeName.String())
	if m.TargetName == nil {
		w.Bool(false)
	} else {
		w.Bool(true)
		w.String(m.TargetName.String())
	}
	return nil
}
