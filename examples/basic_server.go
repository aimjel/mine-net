package main

import (
	"fmt"
	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
)

func main() {
	lc := minecraft.ListenConfig{
		OnlineMode:           true,
		CompressionThreshold: 0, //compresses everything!
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

		fmt.Println("accepting connection!")
		go handleConn(c)
	}
}

func handleConn(c *minecraft.Conn) {
	if err := c.SendPacket(&packet.JoinGame{
		GameMode:       1, //creative
		DimensionNames: []string{"minecraft:the_end"},
		DimensionType:  "minecraft:the_end",
		DimensionName:  "earth:itsgoingblow",
	}); err != nil {
		c.Close(err)
		return
	}

	if err := c.SendPacket(&packet.SetDefaultSpawnPosition{}); err != nil {
		c.Close(err)
		return
	}

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			fmt.Printf("%v: %v", c.RemoteAddr(), err)
			c.Close(nil)
			return
		}

		fmt.Printf("packet id %#v\n", pk.ID())
	}
}
