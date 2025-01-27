package packet

import (
	"github.com/aimjel/minenet/chat"
	"github.com/aimjel/minenet/protocol/encoding"
)

type CommandSuggestionsResponse struct {
	TransactionId int32
	Start         int32
	Length        int32
	Matches       []SuggestionMatch
}

type SuggestionMatch struct {
	Match   string
	Tooltip string
}

func (m CommandSuggestionsResponse) ID() int32 {
	return 0x0F
}

func (m *CommandSuggestionsResponse) Decode(r *encoding.Reader) error {
	// todo implement
	return nil
}

func (m CommandSuggestionsResponse) Encode(w *encoding.Writer) error {
	w.VarInt(m.TransactionId)
	w.VarInt(m.Start)
	w.VarInt(m.Length)
	w.VarInt(int32(len(m.Matches)))
	for _, match := range m.Matches {
		w.String(match.Match)
		w.Bool(match.Tooltip != "")
		if match.Tooltip != "" {
			msg := chat.NewMessage(match.Tooltip)
			w.String(msg.String())
		}
	}
	return nil
}
