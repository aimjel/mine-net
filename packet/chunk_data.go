package packet

import "fmt"

type ChunkData struct {
	X, Z int32

	Heightmaps []byte

	Sections []byte

	//Todo block entities, chests etc
	Blocks struct{}

	Lights []byte
}

func (d ChunkData) ID() int32 {
	return 0x24
}

func (d *ChunkData) Decode(r *Reader) error {
	panic("implement") //todo implement decode chunk data packet
	return nil
}

func (d *ChunkData) Encode(w *Writer) error {
	_ = w.Int32(d.X)
	_ = w.Int32(d.Z)

	_ = w.FixedByteArray(d.Heightmaps)
	fmt.Println(len(d.Heightmaps), "len of height map after writing to packet writer buffer")

	_ = w.ByteArray(d.Sections)

	//TODO block entities
	_ = w.VarInt(0)

	return w.FixedByteArray(d.Lights)
}
