package packet

type PlayerMovement struct {
	OnGround bool
}

func (m PlayerMovement) ID() int32 {
	return 0x14
}

func (m *PlayerMovement) Decode(r *Reader) error {
	return r.Bool(&m.OnGround)
}

func (m PlayerMovement) Encode(w Writer) error {
	return w.Bool(m.OnGround)
}
