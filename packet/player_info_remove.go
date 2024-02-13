package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type PlayerInfoRemove struct {
	UUIDs [][16]byte
}

func (p PlayerInfoRemove) ID() int32 {
	return 0x39
}

func (p PlayerInfoRemove) Decode(r *encoding.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (p PlayerInfoRemove) Encode(w *encoding.Writer) error {
	_ = w.VarInt(int32(len(p.UUIDs)))

	var err error
	for _, v := range p.UUIDs {
		err = w.UUID(v)
	}

	return err
}
