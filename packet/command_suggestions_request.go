package packet

type CommandSuggestionsRequest struct {
	TransactionId int32
	Text          string
}

func (m CommandSuggestionsRequest) ID() int32 {
	return 0x09
}

func (m *CommandSuggestionsRequest) Decode(r *Reader) error {
	r.VarInt(&m.TransactionId)
	return r.String(&m.Text)
}

func (m CommandSuggestionsRequest) Encode(w Writer) error {
	w.VarInt(m.TransactionId)
	return w.String(m.Text)
}
