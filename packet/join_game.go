package packet

import (
	_ "embed"
	"github.com/aimjel/minecraft/protocol/encoding"
)

//go:embed internal/data/dimension.nbt
var dimensions []byte

type JoinGame struct {
	EntityID         int32
	IsHardcore       bool
	GameMode         uint8
	PreviousGameMode int8
	DimensionNames   []string

	Registry []byte

	DimensionType       string
	DimensionName       string
	HashedSeed          int64
	MaxPlayers          int32
	ViewDistance        int32
	SimulationDistance  int32
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	IsDebug             bool
	IsFlat              bool
	DeathDimensionName  string
	DeathLocation       uint64
	PartialCooldown     int32
}

func (g JoinGame) ID() int32 {
	return 0x28
}

func (g *JoinGame) Decode(r *encoding.Reader) error {
	panic("implement") //todo implement decode join game packet
	return nil
}

func (g JoinGame) Encode(w *encoding.Writer) error {
	_ = w.Int32(g.EntityID)
	_ = w.Bool(g.IsHardcore)
	_ = w.Uint8(g.GameMode)
	_ = w.Int8(g.PreviousGameMode)
	_ = w.VarInt(int32(len(g.DimensionNames)))
	for _, world := range g.DimensionNames {
		_ = w.String(world)
	}

	if g.Registry != nil {
		_ = w.FixedByteArray(g.Registry)
	} else {
		_ = w.FixedByteArray(dimensions)
	}

	_ = w.String(g.DimensionType)
	_ = w.String(g.DimensionName)
	_ = w.Int64(g.HashedSeed)
	_ = w.VarInt(g.MaxPlayers)
	_ = w.VarInt(g.ViewDistance)
	_ = w.VarInt(g.SimulationDistance)
	_ = w.Bool(g.ReducedDebugInfo)
	_ = w.Bool(g.EnableRespawnScreen)
	_ = w.Bool(g.IsDebug)
	_ = w.Bool(g.IsFlat)

	if g.DeathDimensionName != "" && g.DeathLocation != 0 {
		_ = w.Bool(true)
		_ = w.String(g.DeathDimensionName)
		_ = w.Uint64(g.DeathLocation)
	} else {
		_ = w.Bool(false)
	}

	return w.VarInt(g.PartialCooldown)
}
