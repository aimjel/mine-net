package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type Handshake struct {
	ProtocolVersion int32

	ServerAddress string

	ServerPort uint16

	NextState uint8
}

func (h Handshake) ID() int32 {
	return 0x00

}

func (h *Handshake) Decode(r *encoding.Reader) error {
	_ = r.VarInt(&h.ProtocolVersion)
	_ = r.String(&h.ServerAddress)
	_ = r.Uint16(&h.ServerPort)
	return r.Uint8(&h.NextState)
}

func (h Handshake) Encode(w *encoding.Writer) error {
	_ = w.VarInt(h.ProtocolVersion)
	_ = w.String(h.ServerAddress)
	_ = w.Uint16(h.ServerPort)
	return w.Uint8(h.NextState)
}
