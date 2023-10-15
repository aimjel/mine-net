package packet

type SpawnEntity struct {
	EntityID                        int32
	UUID                            [16]byte
	Type                            int32
	X, Y, Z                         float64
	Pitch, Yaw, HeadYaw             byte
	Data                            int32
	VelocityX, VelocityY, VelocityZ int16
}

func (e *SpawnEntity) ID() int32 {
	return 0x01
}

func (e *SpawnEntity) Decode(r *Reader) error {
	return nil
}

func (e *SpawnEntity) Encode(w Writer) error {
	_ = w.VarInt(e.EntityID)
	_ = w.UUID(e.UUID)
	_ = w.VarInt(e.Type)
	_ = w.Float64(e.X)
	_ = w.Float64(e.Y)
	_ = w.Float64(e.Z)
	_ = w.Uint8(e.Pitch)
	_ = w.Uint8(e.Yaw)
	_ = w.Uint8(e.HeadYaw)
	_ = w.VarInt(e.Data)
	_ = w.Int16(e.VelocityX)
	_ = w.Int16(e.VelocityY)
	return w.Int16(e.VelocityZ)
}
