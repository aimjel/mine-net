package packet

import (
	"errors"
	"github.com/aimjel/minenet/protocol/encoding"
)

var NotImplemented = errors.New("a packet field has not been implemented")

type Packet interface {
	ID() int32

	Decode(r *encoding.Reader) error

	Encode(w *encoding.Writer) error
}

type Unknown struct {
	Id      int32
	Payload []byte
}

func (u Unknown) ID() int32 {
	return u.Id
}

func (u Unknown) Decode(*encoding.Reader) error {
	return nil
}

func (u Unknown) Encode(w *encoding.Writer) error {
	return w.FixedByteArray(u.Payload)
}
