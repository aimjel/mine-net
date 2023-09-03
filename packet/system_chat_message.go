package packet

import "github.com/aimjel/minecraft/chat"

type SystemChatMessage struct {
	Content string
}

func (m SystemChatMessage) ID() int32 {
	return 0x64
}

func (m *SystemChatMessage) Decode(r *Reader) error {
	return r.String(&m.Content)
}

func (m SystemChatMessage) Encode(w Writer) error {
	content := chat.NewMessage(m.Content)
	w.String(content.String())
	return w.Bool(false)
}
