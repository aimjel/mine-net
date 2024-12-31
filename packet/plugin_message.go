package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type PlayClientboundPluginMessage struct {
	Channel string
	Data    []byte
}

func (p PlayClientboundPluginMessage) ID() int32 {
	return 0x19
}

func (p PlayClientboundPluginMessage) Decode(r *encoding.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (p PlayClientboundPluginMessage) Encode(w *encoding.Writer) error {
	_ = w.String(p.Channel)
	return w.FixedByteArray(p.Data)
}

type ConfigurationClientboundPluginMessage struct {
	Channel string
	Data    []byte
}

func (p ConfigurationClientboundPluginMessage) ID() int32 {
	return 0x01
}

func (p ConfigurationClientboundPluginMessage) Decode(r *encoding.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (p ConfigurationClientboundPluginMessage) Encode(w *encoding.Writer) error {
	_ = w.String(p.Channel)
	return w.FixedByteArray(p.Data)
}
