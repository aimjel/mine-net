package packet

import "github.com/aimjel/minenet/protocol/encoding"

type DamageEvent struct {
	EntityID       int32
	SourceTypeID   int32
	SourceCauseID  int32
	SourceDirectID int32

	HasSrcPos bool
	X, Y, Z   float64
}

func (l DamageEvent) ID() int32 {
	return 0x18
}

func (l *DamageEvent) Decode(r *encoding.Reader) error {
	_ = r.VarInt(&l.EntityID)
	_ = r.VarInt(&l.SourceTypeID)
	_ = r.VarInt(&l.SourceCauseID)
	_ = r.VarInt(&l.SourceDirectID)
	_ = r.Bool(&l.HasSrcPos)
	if l.HasSrcPos {
		_ = r.Float64(&l.X)
		_ = r.Float64(&l.Y)
		_ = r.Float64(&l.Z)
	}
	return nil
}

func (l DamageEvent) Encode(w *encoding.Writer) error {
	_ = w.VarInt(l.EntityID)
	_ = w.VarInt(l.SourceTypeID)
	_ = w.VarInt(l.SourceCauseID)
	_ = w.VarInt(l.SourceDirectID)
	_ = w.Bool(l.HasSrcPos)
	if l.HasSrcPos {
		_ = w.Float64(l.X)
		_ = w.Float64(l.Y)
		_ = w.Float64(l.Z)
	}
	return nil
}
