package packet

type UnloadChunk struct {
	ChunkX, ChunkZ int32
}

func (c UnloadChunk) ID() int32 {
	return 0x1E
}

func (c *UnloadChunk) Decode(r *Reader) error {
	_ = r.VarInt(&c.ChunkX)
	return r.VarInt(&c.ChunkZ)
}

func (c UnloadChunk) Encode(w *Writer) error {
	_ = w.Int32(c.ChunkX)
	return w.Int32(c.ChunkZ)
}
