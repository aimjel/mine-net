package metadata

import (
	"github.com/aimjel/minecraft/chat"
	"github.com/aimjel/minecraft/protocol/encoding"
)

type pose uint8

const (
	standing pose = iota
	fallFlying
	sleeping
	swimming
	spinAttack
	sneaking
	longJumping
	dying
	croaking
	usingTongue
	sitting
	roaring
	sniffing
	emerging
	digging
)

type entityData = uint8

const (
	onFire entityData = 1 << iota
	isCrouching
	unused
	isSprinting
	isSwimming
	isInvisible
	glowingEffect
	flyingWithElytra
)

type entityIndex uint8

const (
	data entityIndex = 1 << iota
	airTicks
	customName
	customNameVisible
	silent
	noGravity
	poseType
	frozenTicks
)

// Entity implements the base entity metadata values
type Entity struct {
	data              entityData
	airTicks          int32
	customName        *chat.Message
	customNameVisible bool
	silent            bool
	noGravity         bool
	pose              pose
	frozenTicks       int32

	indexUsed entityIndex
}

func (e Entity) Encode(w *encoding.Writer) error {
	if e.indexUsed&data != 0 {
		_ = w.Uint8(bitmaskToIndex(data))
		_ = encode(w, e.data)
	}
	if e.indexUsed&airTicks != 0 {
		_ = w.Uint8(bitmaskToIndex(airTicks))
		_ = encode(w, e.airTicks)
	}
	if e.indexUsed&customName != 0 {
		_ = w.Uint8(bitmaskToIndex(customName))
		_ = encode(w, e.customName)
	}
	if e.indexUsed&customNameVisible != 0 {
		_ = w.Uint8(bitmaskToIndex(customNameVisible))
		_ = encode(w, e.customNameVisible)
	}
	if e.indexUsed&silent != 0 {
		_ = w.Uint8(bitmaskToIndex(silent))
		_ = encode(w, e.silent)
	}
	if e.indexUsed&noGravity != 0 {
		_ = w.Uint8(bitmaskToIndex(noGravity))
		_ = encode(w, e.noGravity)
	}
	if e.indexUsed&poseType != 0 {
		_ = w.Uint8(bitmaskToIndex(poseType))
		_ = encode(w, e.pose)
	}
	if e.indexUsed&frozenTicks != 0 {
		_ = w.Uint8(bitmaskToIndex(poseType))
		_ = encode(w, e.frozenTicks)
	}

	return nil
}

func (e Entity) Decode(r *encoding.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (e Entity) Crouch(v bool) Entity {
	if v {
		e.data |= isCrouching
		e.pose = sneaking
	} else {
		//don't need to set data
		e.pose = standing
	}
	e.indexUsed |= poseType | data

	return e
}

func (e Entity) Sprinting(v bool) Entity {
	if v {
		e.data |= isSprinting
	}

	e.indexUsed |= data

	return e
}
