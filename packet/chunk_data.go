package packet

import "github.com/aimjel/minenet/protocol/encoding"

type ChunkData struct {
	X, Z int32

	//Data includes height-map, section data, block entities
	// and light information.
	Data []byte
}

func (d ChunkData) ID() int32 {
	return 0x24
}

func (d *ChunkData) Decode(r *encoding.Reader) error {
	panic("implement") //todo implement decode chunk data packet
	return nil
}

func (d *ChunkData) Encode(w *encoding.Writer) error {
	_ = w.Int32(d.X)
	_ = w.Int32(d.Z)

	return w.FixedByteArray(d.Data)
}
