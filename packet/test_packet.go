package packet

// TestPacket is a packet which holds every possible
// data type which can be received or sent to the client
type TestPacket struct {
	Boolean       bool
	Byte          int8
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16
	Int           int32
	Long          int64
	Float         float32
	Double        float64
	String        string
	//skipping chat since it's encoded like a string
	//? Implement a custom chat data type?
	//skipping identifier since it's encoded like a string
	//? Implement a custom chat data type?

	VarInt  int32
	VarLong int64

	//todo position and angle

	UUID [16]byte

	ByteArray []uint8
}

func (t *TestPacket) ID() int32 {
	return -0x01
}

func (t *TestPacket) Decode(r *Reader) error {
	r.Bool(&t.Boolean)
	r.Int8(&t.Byte)
	r.Uint8(&t.UnsignedByte)
	r.Int16(&t.Short)
	r.Uint16(&t.UnsignedShort)
	r.Int32(&t.Int)
	r.Int64(&t.Long)
	r.Float32(&t.Float)
	r.Float64(&t.Double)
	r.String(&t.String)

	//todo chat and identifier

	r.VarInt(&t.VarInt)

	//todo varlong, position, angle

	r.UUID(&t.UUID)
	return r.ByteArray(&t.ByteArray)
}

func (t *TestPacket) Encode(w *Writer) error {
	w.Bool(t.Boolean)
	w.Int8(t.Byte)
	w.Uint8(t.UnsignedByte)
	w.Int16(t.Short)
	w.Uint16(t.UnsignedShort)
	w.Int32(t.Int)
	w.Int64(t.Long)
	w.Float32(t.Float)
	w.Float64(t.Double)
	w.String(t.String)

	//todo chat and identifier

	w.VarInt(t.VarInt)

	//todo varlong, position, angle

	w.UUID(t.UUID)
	return w.ByteArray(t.ByteArray)
}
