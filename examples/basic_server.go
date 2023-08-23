package main

import (
	"fmt"
	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
)

func main() {
	lc := minecraft.ListenConfig{
		OnlineMode:           true,
		CompressionThreshold: 256, //compresses everything!
		Status:               minecraft.NewStatus(756, 10, "someone had todo it"),
	}

	l, err := lc.Listen("localhost:25565")
	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("main thread|server: %v", err)
			return
		}

		if err := c.SendPacket(&packet.JoinGame{
			DimensionNames: []string{"minecraft:overworld"},
			DimensionName:  "minecraft:overworld",
		}); err != nil {
			c.Close(err)
		}

		if err := c.SendPacket(&packet.PlayerPositionLook{}); err != nil {
			c.Close(err)
		}
	}
}
