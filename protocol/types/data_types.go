package types

type Position int64

func (p Position) XYZ() (int32, int32, int32) {
	return int32(p >> 38), int32(p << 52 >> 52), int32(p << 26 >> 38)
}
