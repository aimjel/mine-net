package packet

import _ "embed"

//go:embed internal/data/dimension.nbt
var dimensions []byte

type JoinGame struct {
	EntityID         int32
	IsHardcore       bool
	GameMode         uint8
	PreviousGameMode int8
	WorldNames       []string

	//Dimension Codec & Dimension

	DimensionName       string
	HashedLong          int64
	MaxPlayers          int32
	ViewDistance        int32
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	IsDebug             bool
	IsFlat              bool
}

func (g JoinGame) ID() int32 {
	return 0x26
}

func (g *JoinGame) Decode(r *Reader) error {
	panic("implement") //todo implement decode join game packet
	return nil
}

func (g JoinGame) Encode(w Writer) error {
	_ = w.Int32(g.EntityID)
	_ = w.Bool(g.IsHardcore)
	_ = w.Uint8(g.GameMode)
	_ = w.Int8(g.PreviousGameMode)
	_ = w.VarInt(int32(len(g.WorldNames)))
	for _, world := range g.WorldNames {
		_ = w.String(world)
	}
	_ = w.Nbt(dimensions)
	_ = w.String(g.DimensionName)
	_ = w.Int64(g.HashedLong)
	_ = w.VarInt(g.MaxPlayers)
	_ = w.VarInt(g.ViewDistance)
	_ = w.Bool(g.IsHardcore)
	_ = w.Bool(g.ReducedDebugInfo)
	_ = w.Bool(g.IsDebug)
	return w.Bool(g.IsFlat)
}
