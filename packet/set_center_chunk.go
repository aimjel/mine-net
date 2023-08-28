package packet

type SetCenterChunk struct {
	ChunkX int32
	ChunkZ int32
}

func (p SetCenterChunk) ID() int32 {
	return 0x4E
}

func (p *SetCenterChunk) Decode(r *Reader) error {
	_ = r.VarInt(&p.ChunkX)
	return r.VarInt(&p.ChunkZ)
}

func (p SetCenterChunk) Encode(w Writer) error {
	_ = w.VarInt(p.ChunkX)
	return w.VarInt(p.ChunkZ)
}
