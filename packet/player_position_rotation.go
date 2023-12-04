package packet

type PlayerPositionRotation struct {
	X, FeetY, Z float64

	Yaw, Pitch float32

	OnGround bool
}

func (r PlayerPositionRotation) ID() int32 {
	return 0x15
}

func (r *PlayerPositionRotation) Decode(rd *Reader) error {
	_ = rd.Float64(&r.X)
	_ = rd.Float64(&r.FeetY)
	_ = rd.Float64(&r.Z)
	_ = rd.Float32(&r.Yaw)
	_ = rd.Float32(&r.Pitch)

	return rd.Bool(&r.OnGround)
}

func (r PlayerPositionRotation) Encode(w *Writer) error {
	_ = w.Float64(r.X)
	_ = w.Float64(r.FeetY)
	_ = w.Float64(r.Z)
	_ = w.Float32(r.Yaw)
	_ = w.Float32(r.Pitch)

	return w.Bool(r.OnGround)
}
