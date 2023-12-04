package packet

type ChangeDifficulty struct {
	Difficulty       uint8
	DifficultyLocked bool
}

func (m ChangeDifficulty) ID() int32 {
	return 0x0C
}

func (m *ChangeDifficulty) Decode(r *Reader) error {
	r.Uint8(&m.Difficulty)
	return r.Bool(&m.DifficultyLocked)
}

func (m ChangeDifficulty) Encode(w *Writer) error {
	w.Uint8(m.Difficulty)
	return w.Bool(m.DifficultyLocked)
}
