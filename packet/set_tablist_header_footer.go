package packet

import "github.com/aimjel/minecraft/chat"

type SetTablistHeaderFooter struct {
	Header string
	Footer string
}

func (m SetTablistHeaderFooter) ID() int32 {
	return 0x65
}

func (m *SetTablistHeaderFooter) Decode(r *Reader) error {
	r.String(&m.Header)
	return r.String(&m.Footer)
}

func (m SetTablistHeaderFooter) Encode(w Writer) error {
	header := chat.NewMessage(m.Header)
	footer := chat.NewMessage(m.Footer)
	w.String(header.String())
	return w.String(footer.String())
}
