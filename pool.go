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
	0x07: func() packet.Packet { return &packet.ClientStatus{} },
	0x08: func() packet.Packet { return &packet.ClientSettings{} },

	0x14: func() packet.Packet { return &packet.PlayerPosition{} },
	0x15: func() packet.Packet { return &packet.PlayerPositionRotation{} },
	0x16: func() packet.Packet { return &packet.PlayerRotation{} },
	0x17: func() packet.Packet { return &packet.PlayerMovement{} },

	0x12: func() packet.Packet { return &packet.KeepAlive{} },

	0x4d: func() packet.Packet { return &packet.HeldItemChange{} },

	0x2b: func() packet.Packet { return &packet.CreateInventoryAction{} },

	0x31: func() packet.Packet { return &packet.PlayerBlockPlacement{} },
	0x32: func() packet.Packet { return &packet.UseItem{} },
}

type clientLoginPool struct{}

func (c clientLoginPool) Get(id int32) packet.Packet {
	return nil
}
