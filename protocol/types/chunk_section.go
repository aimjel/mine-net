package types

type ChunkSection struct {
	BlockCount uint16

	BlockStates PaletteContainer

	Biomes PaletteContainer

	SkyLight []int8

	BlockLight []int8
}

type PaletteContainer struct {
	BitsPerEntry uint8

	//entries global ids
	Entries []int32

	Data []int64
}
