package packet

import "github.com/aimjel/minenet/protocol/encoding"

type ChangeDifficulty struct {
	Difficulty       uint8
	DifficultyLocked bool
}

func (m ChangeDifficulty) ID() int32 {
	return 0x0C
}

func (m *ChangeDifficulty) Decode(r *encoding.Reader) error {
	r.Uint8(&m.Difficulty)
	return r.Bool(&m.DifficultyLocked)
}

func (m ChangeDifficulty) Encode(w *encoding.Writer) error {
	w.Uint8(m.Difficulty)
	return w.Bool(m.DifficultyLocked)
}
