package packet

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

type SystemChatMessage struct {
	Message chat.Message
}

func (m SystemChatMessage) ID() int32 {
	return 0x64
}

func (m *SystemChatMessage) Decode(r *encoding.Reader) error {
	return NotImplemented
}

func (m SystemChatMessage) Encode(w *encoding.Writer) error {
	w.String(m.Message.String())
	return w.Bool(false)
}
