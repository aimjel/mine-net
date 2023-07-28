package protocol

import "github.com/aimjel/minecraft/packet"

type Pool struct {
	packets []packet.Packet

	version int32
}

func NewPool(in []packet.Packet) Pool {
	return Pool{packets: in}
}

func (p *Pool) Get(id int32) packet.Packet {
	if int32(len(p.packets)) < id {
		return nil
	}

	return p.packets[id]
}
