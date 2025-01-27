package metadata

import (
	"github.com/aimjel/minenet/protocol/encoding"
	"math/bits"
)

type MetaData interface {
	Encode(w *encoding.Writer) error
	Decode(r *encoding.Reader) error
}

func bitmaskToIndex[T int | entityIndex](bitmask T) byte {
	if bitmask == 0 {
		return 0 // Handle case when there is no set bit
	}
	return byte(bits.Len(uint(bitmask)) - 1)
}
