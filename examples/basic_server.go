package main

import (
	"fmt"
	"time"

	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
)

func main() {
	lc := minecraft.ListenConfig{
		OnlineMode:           true,
		CompressionThreshold: 256, //compresses everything!
		Status:               minecraft.NewStatus(minecraft.Version{Protocol: 763}, 10, "someone had todo it"),
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
			GameMode:       1, //creative
			DimensionNames: []string{"minecraft:the_end"},
			DimensionType:  "minecraft:the_end",
			DimensionName:  "earth:itsgoingblow",
		}); err != nil {
			c.Close(err)
		}

		if err := c.SendPacket(&packet.SetDefaultSpawnPosition{}); err != nil {
			c.Close(err)
		}

		c.Close(nil)

		time.Sleep(5 * time.Second)

		if err := c.SendPacket(&packet.JoinGame{
			GameMode:       1, //creative
			DimensionNames: []string{"minecraft:overworld"},
			DimensionType:  "minecraft:overworld",
			DimensionName:  "earth:itsgoingblow",
		}); err != nil {
			c.Close(err)
		}

		if err := c.SendPacket(&packet.SetDefaultSpawnPosition{}); err != nil {
			c.Close(err)
		}
	}
}
