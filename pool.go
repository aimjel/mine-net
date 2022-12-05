package minecraft

import (
	"github.com/aimjel/minecraft/packet"
)

var handshakePool = map[int32]func() packet.Packet{
	0x00: func() packet.Packet { return &packet.Handshake{} },
}

var serverBoundPlayPool = map[int32]func() packet.Packet{
	0x03: func() packet.Packet { return &packet.ChatMessageServer{} },
	0x04: func() packet.Packet { return &packet.ClientStatus{} },
	0x05: func() packet.Packet { return &packet.ClientSettings{} },

	0x11: func() packet.Packet { return &packet.PlayerPosition{} },
	0x12: func() packet.Packet { return &packet.PlayerPositionRotation{} },
	0x13: func() packet.Packet { return &packet.PlayerRotation{} },
	0x14: func() packet.Packet { return &packet.PlayerMovement{} },

	0x0F: func() packet.Packet { return &packet.KeepAlive{} },

	0x25: func() packet.Packet { return &packet.HeldItemChange{} },

	0x28: func() packet.Packet { return &packet.CreateInventoryAction{} },

	0x2E: func() packet.Packet { return &packet.PlayerBlockPlacement{} },
	0x2F: func() packet.Packet { return &packet.UseItem{} },
}
