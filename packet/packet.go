package packet

import "errors"

var NotImplemented = errors.New("a packet field has not been implemented")

type Packet interface {
	ID() int32

	Decode(r *Reader) error

	Encode(w *Writer) error
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

func (u Unknown) Encode(w *Writer) error {
	return w.FixedByteArray(u.Payload)
}
