package packet

import (
	"github.com/aimjel/minecraft/protocol/encoding"
	"github.com/aimjel/minecraft/protocol/types"
)

type SetDefaultSpawnPosition struct {
	Location types.Position
	Angle    float32
}

func (s SetDefaultSpawnPosition) ID() int32 {
	return 0x50
}

func (s SetDefaultSpawnPosition) Decode(r *encoding.Reader) error {
	_ = r.Int64((*int64)(&s.Location))
	return r.Float32(&s.Angle)
}

func (s SetDefaultSpawnPosition) Encode(w *encoding.Writer) error {
	_ = w.Int64(int64(s.Location))
	return w.Float32(s.Angle)
}
