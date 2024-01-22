package packet

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

type ServerData struct {
	MOTD               chat.Message
	Icon               []byte
	EnforcesSecureChat bool
}

func (m ServerData) ID() int32 {
	return 0x45
}

func (m *ServerData) Decode(r *encoding.Reader) error {
	return NotImplemented
}

func (m ServerData) Encode(w *encoding.Writer) error {
	w.String(m.MOTD.String())
	if m.Icon == nil {
		w.Bool(false)
	} else {
		w.Bool(true)
		w.ByteArray(m.Icon)
	}
	w.Bool(m.EnforcesSecureChat)
	return nil
}
