package packet

import "github.com/aimjel/minenet/protocol/encoding"

type SetHealth struct {
	Health         float32
	Food           int32
	FoodSaturation float32
}

func (h SetHealth) ID() int32 {
	return 0x57
}

func (h *SetHealth) Decode(r *encoding.Reader) error {
	r.Float32(&h.Health)
	r.VarInt(&h.Food)
	return r.Float32(&h.FoodSaturation)
}

func (h SetHealth) Encode(w *encoding.Writer) error {
	w.Float32(h.Health)
	w.VarInt(h.Food)
	return w.Float32(h.FoodSaturation)
}
