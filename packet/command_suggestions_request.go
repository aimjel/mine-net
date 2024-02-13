package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type CommandSuggestionsRequest struct {
	TransactionId int32
	Text          string
}

func (m CommandSuggestionsRequest) ID() int32 {
	return 0x09
}

func (m *CommandSuggestionsRequest) Decode(r *encoding.Reader) error {
	r.VarInt(&m.TransactionId)
	return r.String(&m.Text)
}

func (m CommandSuggestionsRequest) Encode(w *encoding.Writer) error {
	w.VarInt(m.TransactionId)
	return w.String(m.Text)
}
