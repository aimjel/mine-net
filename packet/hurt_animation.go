package packet

import (
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/chat"
)

type HurtAnimation struct {
	EntityID int32
  Yaw float32
}

func (l HurtAnimation) ID() int32 {
	return 0x21
}

func (l *HurtAnimation) Decode(r *Reader) error {
  r.VarInt(&l.EntityID)
	return r.Float32(&l.Yaw)
}

func (l HurtAnimation) Encode(w Writer) error {
  w.VarInt(l.EntityID)
	return w.Float32(l.Yaw)
}
