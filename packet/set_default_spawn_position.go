package packet

type SetDefaultSpawnPosition struct {
	Location uint64
	Angle    float32
}

func (s SetDefaultSpawnPosition) ID() int32 {
	return 0x50
}

func (s SetDefaultSpawnPosition) Decode(r *Reader) error {
	//TODO implement me
	panic("implement me")
}

func (s SetDefaultSpawnPosition) Encode(w Writer) error {
	_ = w.Uint64(s.Location)
	return w.Float32(s.Angle)
}
