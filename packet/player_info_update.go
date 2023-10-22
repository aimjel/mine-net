package packet

import (
	"github.com/aimjel/minecraft/protocol/types"
)

type PlayerInfoUpdate struct {
	Actions byte
	Players []types.PlayerInfo
}

func (i *PlayerInfoUpdate) ID() int32 {
	return 0x3A
}

func (i *PlayerInfoUpdate) Decode(r *Reader) error {
	panic("implement me")
}

func (i *PlayerInfoUpdate) Encode(w Writer) error {
	_ = w.Uint8(i.Actions)
	_ = w.VarInt(int32(len(i.Players)))
	for _, p := range i.Players {
		_ = w.UUID(p.UUID)

		if i.Actions&0x01 != 0 {
			//add player
			_ = w.String(p.Name)
			_ = w.VarInt(int32(len(p.Properties)))
			for _, v := range p.Properties {
				_ = w.String(v.Name)
				_ = w.String(v.Value)
				_ = w.Bool(v.Signature != "")
				if v.Signature != "" {
					_ = w.String(v.Signature)
				}
			}
		}

		if i.Actions&0x02 != 0 {
			_ = w.Bool(p.KeySignature != nil)
			_ = w.UUID(p.ChatSessionID)
			_ = w.Uint64(uint64(p.ExpiresAt))
			_ = w.ByteArray(p.PublicKey)
			_ = w.ByteArray(p.KeySignature)
		}

		if i.Actions&0x04 != 0 {
			//update game-mode
			_ = w.VarInt(p.GameMode)
		}

		if i.Actions&0x08 != 0 {
			//enables/disables the player on the player list
			_ = w.Bool(p.Listed)
		}

		if i.Actions&0x10 != 0 {
			//updates ping icon
			_ = w.VarInt(p.Ping)
		}

		if i.Actions&0x20 != 0 {
			//updates display name
			_ = w.Bool(p.HasDisplayName)
			if p.HasDisplayName {
				_ = w.String(p.DisplayName)
			}
		}
	}
	return nil
}
