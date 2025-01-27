package packet

import (
	"github.com/aimjel/minenet/chat"
	"github.com/aimjel/minenet/protocol/encoding"
)

type CombatDeath struct {
	PlayerID int32
	Message  string
}

func (c CombatDeath) ID() int32 {
	return 0x38
}

func (c *CombatDeath) Decode(r *encoding.Reader) error {
	return r.VarInt(&c.PlayerID)
	// todo implement message
}

func (c CombatDeath) Encode(w *encoding.Writer) error {
	w.VarInt(c.PlayerID)
	msg := chat.NewMessage(c.Message)
	return w.String(msg.String())
}
