package packet

import "github.com/aimjel/minenet/protocol/encoding"

type ResourcePackResult struct {
	Result int32
}

func (m ResourcePackResult) ID() int32 {
	return 0x24
}

func (m *ResourcePackResult) Decode(r *encoding.Reader) error {
	return r.VarInt(&m.Result)
}

func (m ResourcePackResult) Encode(w *encoding.Writer) error {
	return w.VarInt(m.Result)
}
