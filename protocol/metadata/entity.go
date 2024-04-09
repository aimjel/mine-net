package metadata

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

type pose = int32

const (
	Standing pose = iota
	FallFlying
	Sleeping
	Swimming
	SpinAttack
	Sneaking
	LongJumping
	Dying
	Croaking
	UsingTongue
	Sitting
	Roaring
	Sniffing
	Emerging
	Digging
)

type EntityDataType = uint8

const (
	OnFire EntityDataType = 1 << iota
	IsCrouching
	unused
	IsSprinting
	IsSwimming
	IsInvisible
	GlowingEffect
	FlyingWithElytra
)

type entityIndex uint8

const (
	EntityData entityIndex = 1 << iota
	AirTicks
	CustomName
	CustomNameVisible
	Silent
	NoGravity
	PoseType
	FrozenTicks
)

// Entity implements the base entity metadata values
type Entity struct {
	Data              EntityDataType
	AirTicks          int32
	CustomName        *chat.Message
	CustomNameVisible bool
	Silent            bool
	NoGravity         bool
	Pose              pose
	FrozenTicks       int32

	IndexUsed entityIndex
}

func (e Entity) Encode(w *encoding.Writer) error {
	if e.IndexUsed&EntityData != 0 {
		_ = w.Uint8(bitmaskToIndex(EntityData))
		_ = encode(w, e.Data)
	}
	if e.IndexUsed&AirTicks != 0 {
		_ = w.Uint8(bitmaskToIndex(AirTicks))
		_ = encode(w, e.AirTicks)
	}
	if e.IndexUsed&CustomName != 0 {
		_ = w.Uint8(bitmaskToIndex(CustomName))
		_ = encode(w, e.CustomName)
	}
	if e.IndexUsed&CustomNameVisible != 0 {
		_ = w.Uint8(bitmaskToIndex(CustomNameVisible))
		_ = encode(w, e.CustomNameVisible)
	}
	if e.IndexUsed&Silent != 0 {
		_ = w.Uint8(bitmaskToIndex(Silent))
		_ = encode(w, e.Silent)
	}
	if e.IndexUsed&NoGravity != 0 {
		_ = w.Uint8(bitmaskToIndex(NoGravity))
		_ = encode(w, e.NoGravity)
	}
	if e.indexUsed&PoseType != 0 {
		_ = w.Uint8(bitmaskToIndex(PoseType))
		_ = encode(w, e.Pose)
	}
	if e.IndexUsed&FrozenTicks != 0 {
		_ = w.Uint8(bitmaskToIndex(FrozenTicks))
		_ = encode(w, e.FrozenTicks)
	}

	return nil
}

func (e Entity) Decode(r *encoding.Reader) error {
	//TODO implement me
	panic("implement me")
}
