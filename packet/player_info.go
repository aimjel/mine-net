package packet

import (
	"github.com/aimjel/minecraft/protocol"
)

type PlayerInfo struct {
	Action  int32
	Players []protocol.PlayerInfo
}

func (i *PlayerInfo) ID() int32 {
	return 0x36
}

func (i *PlayerInfo) Decode(r *Reader) error {
	_ = r.VarInt(&i.Action)
	var length int32
	if err := r.VarInt(&length); err != nil {
		return err
	}

	i.Players = make([]protocol.PlayerInfo, length)

	switch i.Action {

	//add player
	case 0:
		for j := int32(0); j < length; j++ {
			p := &i.Players[j]
			_ = r.UUID(&p.UUID)
			_ = r.String(&p.Name)

			var propertiesLength int32
			if err := r.VarInt(&propertiesLength); err != nil {
				return err
			}

			p.Properties = make([]struct {
				Name      string
				Value     string
				Signature string
			}, propertiesLength)

			for k := int32(0); k < propertiesLength; k++ {
				property := &p.Properties[k]
				_ = r.String(&property.Name)
				_ = r.String(&property.Value)

				var isSigned bool
				if err := r.Bool(&isSigned); err != nil {
					return err
				}

				if isSigned {
					_ = r.String(&property.Signature)
				}
			}

			_ = r.VarInt(&p.GameMode)
			_ = r.VarInt(&p.Ping)
			_ = r.Bool(&p.HasDisplayName)
			if p.HasDisplayName {
				_ = r.String(&p.DisplayName)
			}
		}

		//update game-mode
	case 1:
		for _, p := range i.Players {
			_ = r.UUID(&p.UUID)
			_ = r.VarInt(&p.GameMode)
		}

		//update latency
	case 2:
		for _, p := range i.Players {
			_ = r.UUID(&p.UUID)
			_ = r.VarInt(&p.Ping)
		}

		//update display name
	case 3:
		for _, p := range i.Players {
			_ = r.UUID(&p.UUID)
			_ = r.Bool(&p.HasDisplayName)
			if p.HasDisplayName {
				_ = r.String(&p.DisplayName)
			}
		}
	}

	return nil
}

func (i *PlayerInfo) Encode(w Writer) error {
	_ = w.VarInt(i.Action)
	_ = w.VarInt(int32(len(i.Players)))

	switch i.Action {

	//Add player
	case 0:
		for _, p := range i.Players {
			_ = w.UUID(p.UUID)
			_ = w.String(p.Name)

			_ = w.VarInt(int32(len(p.Properties)))

			for _, property := range p.Properties {
				_ = w.String(property.Name)
				_ = w.String(property.Value)
				_ = w.Bool(property.Signature != "")

				if property.Signature != "" {
					_ = w.String(property.Signature)
				}
			}

			_ = w.VarInt(p.GameMode)
			_ = w.VarInt(p.Ping)
			_ = w.Bool(p.HasDisplayName)
			if p.HasDisplayName {
				//TODO: Implement display name field
			}
		}

	//Remove player
	case 4:
		for _, p := range i.Players {
			_ = w.UUID(p.UUID)
		}
	}

	return nil
}
