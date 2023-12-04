package packet

type ClientSettings struct {
	Locale               string
	ViewDistance         int8
	ChatMode             int32
	ChatColors           bool
	DisplayedSkinParts   uint8
	MainHand             int32
	DisableTextFiltering bool
	AllowServerListings  bool
}

func (s ClientSettings) ID() int32 {
	return 0x08
}

func (s *ClientSettings) Decode(r *Reader) error {
	_ = r.String(&s.Locale)
	_ = r.Int8(&s.ViewDistance)
	_ = r.VarInt(&s.ChatMode)
	_ = r.Bool(&s.ChatColors)
	_ = r.Uint8(&s.DisplayedSkinParts)
	_ = r.VarInt(&s.MainHand)
	_ = r.Bool(&s.DisableTextFiltering)

	return r.Bool(&s.AllowServerListings)
}

func (s ClientSettings) Encode(w *Writer) error {
	_ = w.String(s.Locale)
	_ = w.Int8(s.ViewDistance)
	_ = w.VarInt(s.ChatMode)
	_ = w.Bool(s.ChatColors)
	_ = w.Uint8(s.DisplayedSkinParts)
	_ = w.VarInt(s.MainHand)
	_ = w.Bool(s.DisableTextFiltering)
	return w.Bool(s.AllowServerListings)
}
