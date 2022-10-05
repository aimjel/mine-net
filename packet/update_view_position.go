package packet

type UpdateViewPosition struct {
	ChunkX int32
	ChunkZ int32
}

func (p UpdateViewPosition) ID() int32 {
	return 0x49
}

func (p *UpdateViewPosition) Decode(r *Reader) error {
	_ = r.VarInt(&p.ChunkX)
	return r.VarInt(&p.ChunkZ)
}

func (p UpdateViewPosition) Encode(w Writer) error {
	_ = w.VarInt(p.ChunkX)
	return w.VarInt(p.ChunkZ)
}
