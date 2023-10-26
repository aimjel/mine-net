package packet

import "github.com/aimjel/minecraft/chat"

type SystemChatMessage struct {
	Message Message
}

func (m SystemChatMessage) ID() int32 {
	return 0x64
}

func (m *SystemChatMessage) Decode(r *Reader) error {
	return NotImplemented
}

func (m SystemChatMessage) Encode(w Writer) error {
	w.String(m.Message.String())
	return w.Bool(false)
}
