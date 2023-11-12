package packet

type SetDefaultSpawnPosition struct {
	Location int64
	Angle    float32
}

func (s SetDefaultSpawnPosition) ID() int32 {
	return 0x50
}

func (s SetDefaultSpawnPosition) Decode(r *Reader) error {
	_ = r.Int64(&s.Location)
	return r.Float32(&s.Angle)
}

func (s SetDefaultSpawnPosition) Encode(w Writer) error {
	_ = w.Int64(s.Location)
	return w.Float32(s.Angle)
}
