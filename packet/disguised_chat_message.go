package packet

import "github.com/aimjel/minecraft/chat"

type DisguisedChatMessage struct {
	Message string
  ChatType int32
  ChatTypeName string
  TargetName string
}

func (m DisguisedChatMessage) ID() int32 {
	return 0x1B
}

func (m *DisguisedChatMessage) Decode(r *Reader) error {
	return r.String(&m.Message)
}

func (m DisguisedChatMessage) Encode(w Writer) error {
	content := chat.NewMessage(m.Message)
	w.String(content.String())
	w.VarInt(m.ChatType)
  w.String(m.ChatTypeName)
  if m.TargetName != "" {
    w.Bool(false)
  } else {
    w.Bool(true)
    w.String(m.TargetName)
  }
	return nil
}
