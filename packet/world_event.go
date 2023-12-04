package packet

import "github.com/aimjel/minecraft/protocol/types"

type WorldEvent struct {
	Event                 int32
	Location              types.Position
	Data                  int32
	DisableRelativeVolume bool
}

func (c WorldEvent) ID() int32 {
	return 0x25
}

func (c *WorldEvent) Decode(r *Reader) error {
	r.Int32(&c.Event)
	r.Int64((*int64)(&c.Location))
	r.Int32(&c.Data)
	return r.Bool(&c.DisableRelativeVolume)
}

func (c WorldEvent) Encode(w *Writer) error {
	w.Int32(c.Event)
	w.Int64(int64(c.Location))
	w.Int32(c.Data)
	return w.Bool(c.DisableRelativeVolume)
}
