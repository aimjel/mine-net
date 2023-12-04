package packet

type DamageEvent struct {
	EntityID        int32
	SourceTypeID    int32
	SourceCauseID   int32
	SourceDirectID  int32
	SourcePositionX *float64
	SourcePositionY *float64
	SourcePositionZ *float64
}

func (l DamageEvent) ID() int32 {
	return 0x18
}

func (l *DamageEvent) Decode(r *Reader) error {
	return nil
}

func (l DamageEvent) Encode(w *Writer) error {
	w.VarInt(l.EntityID)
	w.VarInt(l.SourceTypeID)
	w.VarInt(l.SourceCauseID)
	w.VarInt(l.SourceDirectID)
	if l.SourcePositionX != nil && l.SourcePositionY != nil && l.SourcePositionZ != nil {
		w.Bool(true)
		w.Float64(*l.SourcePositionX)
		w.Float64(*l.SourcePositionY)
		w.Float64(*l.SourcePositionZ)
	} else {
		w.Bool(false)
	}
	return nil
}
