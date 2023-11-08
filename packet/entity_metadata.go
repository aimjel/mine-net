package packet

import "github.com/aimjel/minecraft/chat"

type SetEntityMetadata struct {
	EntityID           int32
	Pose               *int32
	Data               *byte
	Health             *float32
	DisplayedSkinParts *uint8
	MainHand           *int32
	Slot               *Slot
	HandState          *int8
	CustomName	   *chat.Message
	IsCustomNameVisibile *bool
}

func (*SetEntityMetadata) ID() int32 {
	return 0x52
}

func (*SetEntityMetadata) Decode(*Reader) error {
	return nil
}

func (s SetEntityMetadata) Encode(w Writer) error {
	w.VarInt(s.EntityID)
	if s.Pose != nil {
		w.Uint8(6)
		w.VarInt(20)
		w.VarInt(*s.Pose)
	}
	if s.Data != nil {
		w.Uint8(0)
		w.Uint8(0)
		w.Uint8(*s.Data)
	}
	if s.Health != nil {
		w.Uint8(9)
		w.VarInt(1)
		w.Float32(*s.Health)
	}
	if s.DisplayedSkinParts != nil {
		w.Uint8(17)
		w.VarInt(0)
		w.Uint8(*s.DisplayedSkinParts)
	}
	if s.MainHand != nil {
		w.Uint8(18)
		w.VarInt(0)
		w.Uint8(uint8(*s.MainHand))
	}
	if s.Slot != nil {
			w.Uint8(8)
			w.Uint8(7)
			w.Bool(true)
			w.VarInt(s.Slot.Id)
			w.Int8(s.Slot.Count)
			w.Nbt2(s.Slot.Tag)
	}
	if s.HandState != nil {
		w.Uint8(8)
		w.Uint8(0)
		w.Int8(*s.HandState)
	}
	if s.CustomName != nil {
		w.Uint8(2)
		w.VarInt(6)
		w.Bool(true)
		w.String(s.CustomName.String())
	}
	if s.IsCustomNameVisible != nil {
		w.Uint8(3)
		w.VarInt(8)
		w.Bool(*s.IsCustomNameVisible)
	}
	return w.Uint8(0xFF)
}
