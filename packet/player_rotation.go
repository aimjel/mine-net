package packet

type PlayerRotation struct {
	Yaw, Pitch float32

	OnGround bool
}

func (r PlayerRotation) ID() int32 {
	return 0x16
}

func (r *PlayerRotation) Decode(rd *Reader) error {
	_ = rd.Float32(&r.Yaw)
	_ = rd.Float32(&r.Pitch)

	return rd.Bool(&r.OnGround)
}

func (r PlayerRotation) Encode(w Writer) error {
	_ = w.Float32(r.Yaw)
	_ = w.Float32(r.Pitch)

	return w.Bool(r.OnGround)
}
