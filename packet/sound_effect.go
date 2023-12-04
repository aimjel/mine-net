package packet

type SoundEffect struct {
	SoundId       int32
	SoundCategory int32
	X, Y, Z       int32
	Volume, Pitch float32
}

func (e *SoundEffect) ID() int32 {
	return 0x5c
}

func (e *SoundEffect) Decode(r *Reader) error {
	return nil
}

func (e SoundEffect) Encode(w *Writer) error {
	_ = w.VarInt(e.SoundId)
	_ = w.VarInt(e.SoundCategory)
	_ = w.Int32(e.X)
	_ = w.Int32(e.Y)
	_ = w.Int32(e.Z)
	_ = w.Float32(e.Volume)
	return w.Float32(e.Pitch)
}
