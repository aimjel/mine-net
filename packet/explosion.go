package packet

type Explosion struct {
	X, Y, Z       float32
	Strength      float32
	Records       []byte
	PlayerMotionX float32
	PlayerMotionY float32
	PlayerMotionZ float32
}

func (e *Explosion) ID() int32 {
	return 0x1C
}

func (e *Explosion) Decode(r *Reader) error {
	return nil
}
func (e *Explosion) Encode(w Writer) error {
	_ = w.Float32(e.X)
	_ = w.Float32(e.Y)
	_ = w.Float32(e.Z)
	_ = w.Float32(e.Strength)
	_ = w.ByteArray(e.Records)
	_ = w.Float32(e.PlayerMotionX)
	_ = w.Float32(e.PlayerMotionY)
	return w.Float32(e.PlayerMotionZ)
}
