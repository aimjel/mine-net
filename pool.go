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
	0x05: func() packet.Packet { return &packet.ChatMessageServer{} },
	0x04: func() packet.Packet { return &packet.ClientStatus{} },
	0x08: func() packet.Packet { return &packet.ClientSettings{} },

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

type clientLoginPool struct{}

func (c clientLoginPool) Get(id int32) packet.Packet {
	return nil
}
