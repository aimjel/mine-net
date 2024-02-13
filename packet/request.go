package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type Request struct{}

func (r Request) ID() int32 {
	return 0x00
}

func (r Request) Encode(w *encoding.Writer) error {
	return nil
}

func (r Request) Decode(rd *encoding.Reader) error {
	return nil
}
