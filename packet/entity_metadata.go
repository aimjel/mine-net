package packet

import (
	"github.com/aimjel/minecraft/protocol/encoding"
	"github.com/aimjel/minecraft/protocol/metadata"
)

type SetEntityMetadata struct {
	EntityID int32
	MetaData metadata.MetaData
}

func (*SetEntityMetadata) ID() int32 {
	return 0x52
}

func (s *SetEntityMetadata) Decode(r *encoding.Reader) error {
	_ = r.VarInt(&s.EntityID)
	return s.MetaData.Decode(r)
}

func (s SetEntityMetadata) Encode(w *encoding.Writer) error {
	_ = w.VarInt(s.EntityID)
	_ = s.MetaData.Encode(w)
	return w.Uint8(0xFF)
}
