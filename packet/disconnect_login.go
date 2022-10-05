package packet

import (
	"encoding/json"
	"fmt"
	"minecraft/chat"
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
	b, err := json.Marshal(chat.NewMessage(l.Reason))
	if err != nil {
		return fmt.Errorf("%v marshaling disconnect(login) reason", err)
	}

	return w.ByteArray(b)
}
