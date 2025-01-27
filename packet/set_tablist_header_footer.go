package packet

import (
	"github.com/aimjel/minenet/chat"
	"github.com/aimjel/minenet/protocol/encoding"
)

type SetTablistHeaderFooter struct {
	Header string
	Footer string
}

func (m SetTablistHeaderFooter) ID() int32 {
	return 0x65
}

func (m *SetTablistHeaderFooter) Decode(r *encoding.Reader) error {
	r.String(&m.Header)
	return r.String(&m.Footer)
}

func (m SetTablistHeaderFooter) Encode(w *encoding.Writer) error {
	header := chat.NewMessage(m.Header)
	footer := chat.NewMessage(m.Footer)
	w.String(header.String())
	return w.String(footer.String())
}
