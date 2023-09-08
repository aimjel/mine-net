package packet

type PluginMessage struct {
	Channel string
	Data    []byte
}

func (p PluginMessage) ID() int32 {
	return 0x17
}

func (p PluginMessage) Decode(r *Reader) error {
	//TODO implement me
	panic("implement me")
}

func (p PluginMessage) Encode(w Writer) error {
	_ = w.String(p.Channel)
	return w.ByteArray(p.Data)
}
