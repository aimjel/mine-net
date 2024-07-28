package packet

import (
	"github.com/aimjel/minecraft/protocol/encoding"
	"github.com/aimjel/minecraft/protocol/types"
)

type LoginSuccess struct {
	UUID [16]byte
	Name string

	Properties []types.Property
}

func (s LoginSuccess) ID() int32 {
	return 0x02
}

func (s *LoginSuccess) Decode(r *encoding.Reader) error {
	_ = r.UUID(&s.UUID)
	_ = r.String(&s.Name)

	var length int32
	_ = r.VarInt(&length)
	prpty := make([]types.Property, length)

	for i := int32(0); i < length; i++ {
		p := prpty[i]
		_ = r.String(&p.Name)
		_ = r.String(&p.Value)

		var signed bool
		_ = r.Bool(&signed)
		if signed {
			_ = r.String(&p.Signature)
		}
	}

	s.Properties = prpty
	return nil
}

func (s LoginSuccess) Encode(w *encoding.Writer) error {
	_ = w.UUID(s.UUID)
	_ = w.String(s.Name)

	_ = w.VarInt(int32(len(s.Properties)))

	for i := 0; i < len(s.Properties); i++ {
		p := s.Properties[i]

		_ = w.String(p.Name)
		_ = w.String(p.Value)

		_ = w.Bool(p.Signature != "")
		if p.Signature != "" {
			_ = w.String(p.Signature)
		}

	}
	return nil
}

type LoginSuccess121 struct {
	*LoginSuccess

	StrictErrorHandling bool
}

func (l *LoginSuccess121) Encode(w *encoding.Writer) error {
	if err := l.LoginSuccess.Encode(w); err != nil {
		return err
	}

	return w.Bool(l.StrictErrorHandling)
}

func (l *LoginSuccess121) Decode(r *encoding.Reader) error {
	if err := l.LoginSuccess.Decode(r); err != nil {
		return err
	}

	return r.Bool(&l.StrictErrorHandling)
}
