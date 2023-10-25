package packet

import "errors"

var NotImplemented = errors.New("a packet field has not been implemented")

type Packet interface {
	ID() int32

	Decode(r *Reader) error

	Encode(w Writer) error
}

// calculateVarInts returns the number of bytes the var-int array will use
func sizeVarInts(x []int32) (n int32) {
	for i := 0; i < len(x); i++ {
		n += sizeVarInt(x[i])
	}

	return n
}

func sizeVarInt(x int32) (n int32) {
	ux := uint32(x)
	for ux >= 0x80 {
		n++
		ux >>= 7
	}
	n++
	return
}

type Unknown struct {
	Id      int32
	Payload []byte
}

func (u Unknown) ID() int32 {
	return u.Id
}

func (u Unknown) Decode(*Reader) error {
	return nil
}

func (u Unknown) Encode(w Writer) error {
	return w.FixedByteArray(u.Payload)
}
