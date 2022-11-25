package packet

import "errors"

var NotImplemneted = errors.New("a packet field has not been implemented")

type Packet interface {
	ID() int32

	Decode(r *Reader) error

	Encode(w Writer) error
}

// calculateVarIntLength returns the number of bytes the var-int array will use
func calculateVarIntLength(x []int32) (n int) {
	for i := 0; i < len(x); i++ {

		ux := uint32(x[i])
		for ux >= 0x80 {
			n++
			ux >>= 7
		}
		n++
	}

	return n
}
