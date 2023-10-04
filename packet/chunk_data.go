package packet

import (
	"unsafe"
)

type ChunkData struct {
	X, Z int32

	Heightmaps struct {
		MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
		WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
	}

	Sections []Section

	BlockEntities []struct {
		PackedXZ byte
		Y        uint16
		Type     int32
		Data     []byte //nbt
	}
}

func (d ChunkData) ID() int32 {
	return 0x24
}

func (d *ChunkData) Decode(r *Reader) error {
	panic("implement") //todo implement decode chunk data packet
	return nil
}

func (d ChunkData) Encode(w Writer) error {
	_ = w.Int32(d.X)
	_ = w.Int32(d.Z)

	_ = w.Nbt2(d.Heightmaps)

	var bytes int32
	for _, s := range d.Sections {
		bytes += 2 //block count

		//!BLOCK STATES
		bytes++ //bits per entry

		//palette
		if s.BlockStates.BitsPerEntry == 0 {
			bytes += sizeVarInt(s.BlockStates.Entries[0])
		} else {
			bytes += sizeVarInt(int32(len(s.BlockStates.Entries)))
			bytes += sizeVarInts(s.BlockStates.Entries)
		}

		//data
		if s.BlockStates.BitsPerEntry != 0 {
			bytes += sizeVarInt(int32(len(s.BlockStates.Data)))
			bytes += int32(len(s.BlockStates.Data)) * 8
		} else {
			bytes++ //empty array
		}

		////!BIOMES
		bytes++ //bits per entry

		//palette
		if s.Biomes.BitsPerEntry == 0 {
			//todo fix this
			bytes += sizeVarInt(0)
		} else {
			bytes += sizeVarInt(int32(len(s.Biomes.Entries)))
			bytes += sizeVarInts(s.Biomes.Entries)
		}

		//data
		if s.Biomes.BitsPerEntry != 0 {
			bytes += sizeVarInt(int32(len(s.Biomes.Data)))
			bytes += (int32(s.Biomes.BitsPerEntry) * 64) / 8
		} else {
			bytes++ //empty long
		}
	}

	if err := w.VarInt(bytes); err != nil {
		return err
	}

	for _, s := range d.Sections {
		if s.BlockStates.BitsPerEntry == 0 {
			if s.BlockStates.Entries[0] == 0 {
				_ = w.Uint16(0) //TODO: Implement proper block count
			} else {
				_ = w.Uint16(5000) //TODO: Implement proper block count
			}

			_ = w.Uint8(s.BlockStates.BitsPerEntry)
			_ = w.VarInt(s.BlockStates.Entries[0])
			w.Uint8(0) //empty data array
		} else {
			_ = w.Uint16(5000) //TODO: Implement proper block count

			_ = w.Uint8(s.BlockStates.BitsPerEntry)
			_ = w.VarIntArray(s.BlockStates.Entries)
			_ = w.Int64Array(s.BlockStates.Data)
		}

		_ = w.Uint8(s.Biomes.BitsPerEntry)
		if s.Biomes.BitsPerEntry == 0 {
			_ = w.VarInt(0x39) //biome id
			_ = w.Uint8(0)     //empty data array
		} else {
			_ = w.VarIntArray(s.Biomes.Entries)
			_ = w.Int64Array(s.Biomes.Data)
		}
	}

	//TODO length of block entities
	_ = w.VarInt(0)

	var skyArrays, blockArrays uint8
	skyLight := bitSet{out: make([]int64, 1)}
	blockLight := bitSet{out: make([]int64, 1)}
	emptySkyLight := bitSet{out: make([]int64, 1)}
	emptyBlockLight := bitSet{out: make([]int64, 1)}

	emptySkyLight.set(0)
	emptyBlockLight.set(0)

	emptyBlockLight.set(9) //temp for now
	for i, section := range d.Sections {
		if section.SkyLight != nil {
			skyArrays++
			skyLight.set(i + 1)
		}

		if section.BlockLight != nil {
			blockArrays++
			blockLight.set(i + 1)
		}

		if allZero(section.SkyLight) {
			if section.BlockStates.Entries[0] != 0x00 { //air
				emptySkyLight.set(i + 1)
			}
		}

		if allZero(section.BlockLight) {
			if section.BlockStates.BitsPerEntry != 0 {
				emptyBlockLight.set(i + 1)
			}
		}
	}

	_ = w.Int64Array(skyLight.out)
	_ = w.Int64Array(blockLight.out)
	_ = w.Int64Array(emptySkyLight.out)
	_ = w.Int64Array(emptyBlockLight.out)

	_ = w.Uint8(skyArrays)
	for _, section := range d.Sections {
		if section.SkyLight != nil {
			_ = w.ByteArray(*(*[]byte)(unsafe.Pointer(&section.SkyLight)))
		}
	}

	_ = w.Uint8(blockArrays)
	for _, section := range d.Sections {
		if section.BlockLight != nil {
			_ = w.ByteArray(*(*[]byte)(unsafe.Pointer(&section.BlockLight)))
		}
	}
	return nil
}

type Section struct {
	BlockCount uint16

	BlockStates struct {
		BitsPerEntry uint8
		//entries global ids
		Entries []int32

		Data []int64
	}

	Biomes struct {
		BitsPerEntry uint8
		//entries global ids
		Entries []int32

		Data []int64
	}

	SkyLight []int8

	BlockLight []int8
}

type bitSet struct {
	//the position of the bit we are writing at
	at int64

	out []int64

	//the index of the slice entry we edit
	i int
}

// add a bit to the bit set.
// x parameter controls if the bit should be set
func (b *bitSet) set(x int) {
	if b.at == 64 {
		b.out = append(b.out, 0)
		b.i++
	}

	b.out[b.i] |= 1 << (x % 64)
	b.at++
}

func allZero(s []int8) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}
