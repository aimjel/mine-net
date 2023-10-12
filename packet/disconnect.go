package packet

import (
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/chat"
)

type DisconnectLogin struct {
	Reason string
}

func (l DisconnectLogin) ID() int32 {
	return 0x00
}

func (l *DisconnectLogin) Decode(r *Reader) error {
	return r.String(&l.Reason)
}

func (l DisconnectLogin) Encode(w Writer) error {
	msg := chat.NewMessage(l.Reason)

	return w.String(b.String())
}

type DisconnectPlay struct {
	DisconnectLogin
}

func (p DisconnectPlay) ID() int32 {
	return 0x1a
}
