package packet

type ChunkData struct {
	X, Z int32

	Heightmaps struct {
		Heightmaps struct {
			MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
			WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
		}
	}

	Biomes []int32

	Sections []struct {
		BitsPerEntry uint8
		//Entries global ids
		Entries []int32

		BlockStates []int64
	}

	BlockEntities []byte
}

func (d ChunkData) ID() int32 {
	return 0x22
}

func (d *ChunkData) Decode(r *Reader) error {
	panic("implement") //todo implement decode chunk data packet
	return nil
}

func (d ChunkData) Encode(w Writer) error {
	_ = w.Int32(d.X)
	_ = w.Int32(d.Z)

	//The max height for 1.17.1 is 320.
	//320 / 16 equals only 20 possible sections, so the bit mask will never be over 64 bits(for now)
	_ = w.Uint8(1)
	_ = w.Int64(int64(1<<len(d.Sections) - 1))

	_ = w.Nbt2(d.Heightmaps)

	_ = w.VarIntArray(d.Biomes)

	var dataLength int32
	for _, s := range d.Sections {
		dataLength += 2
		dataLength++

		lengthEntries := int32(calculateVarIntLength(s.Entries))

		dataLength += int32(calculateVarIntLength([]int32{lengthEntries}))
		dataLength += lengthEntries

		lengthStates := int32(len(s.BlockStates) * 8)

		dataLength += int32(calculateVarIntLength([]int32{lengthStates}))

		dataLength += lengthStates
	}

	if err := w.VarInt(dataLength); err != nil {
		return err
	}

	for _, s := range d.Sections {
		_ = w.Uint16(5000) //TODO: Implement proper block count

		_ = w.Uint8(s.BitsPerEntry)

		_ = w.VarIntArray(s.Entries)

		_ = w.Int64Array(s.BlockStates)
	}

	return w.ByteArray(d.BlockEntities)
}
