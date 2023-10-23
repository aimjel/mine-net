package minecraft

import (
	"github.com/aimjel/minecraft/packet"
)

type Pool interface {
	Get(id int32) packet.Packet
}

type basicPool struct{}

func (b basicPool) Get(id int32) packet.Packet {
	if fn, ok := serverBoundPlayPool[id]; ok {
		return fn()
	}

	return nil
}

var serverBoundPlayPool = map[int32]func() packet.Packet{
	0x03: func() packet.Packet { return &packet.MessageAcknowledgment{} },
	0x04: func() packet.Packet { return &packet.ChatCommandServer{} },
	0x05: func() packet.Packet { return &packet.ChatMessageServer{} },
	0x06: func() packet.Packet { return &packet.PlayerSessionServer{} },
	0x07: func() packet.Packet { return &packet.ClientCommandServer{} },
	0x08: func() packet.Packet { return &packet.ClientSettings{} },
	0x09: func() packet.Packet { return &packet.CommandSuggestionsRequest{} },
	0x10: func() packet.Packet { return &packet.InteractServer{} },

	0x12: func() packet.Packet { return &packet.KeepAlive{} },
	0x14: func() packet.Packet { return &packet.PlayerPosition{} },
	0x15: func() packet.Packet { return &packet.PlayerPositionRotation{} },
	0x16: func() packet.Packet { return &packet.PlayerRotation{} },
	0x17: func() packet.Packet { return &packet.PlayerMovement{} },

	0x1C: func() packet.Packet { return &packet.PlayerAbilitiesServer{} },
	0x1D: func() packet.Packet { return &packet.PlayerActionServer{} },
	0x1E: func() packet.Packet { return &packet.PlayerCommandServer{} },

	0x24: func() packet.Packet { return &packet.ResourcePackResult{} },
	0x28: func() packet.Packet { return &packet.SetHeldItemServer{} },
	0x2b: func() packet.Packet { return &packet.SetCreativeModeSlot{} },

	0x2F: func() packet.Packet { return &packet.SwingArmServer{} },

	0x30: func() packet.Packet { return &packet.TeleportToEntityServer{} },

	0x31: func() packet.Packet { return &packet.PlayerBlockPlacement{} },
	0x32: func() packet.Packet { return &packet.UseItem{} },

	0x4d: func() packet.Packet { return &packet.HeldItemChange{} },
}

type nopPool struct{}

func (c nopPool) Get(id int32) packet.Packet {
	return nil
}
