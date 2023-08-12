# Minecraft Protocol Library for Go

Efficient Minecraft protocol integration in Go. 
Seamlessly send, receive, and manipulate packets with its built-in connection class. 
The integrated listener effortlessly manages status and login states, complete with encryption for secure communication and compression for enhanced performance.

## Installation

```sh
go get github.com/aimjel/minecraft
```

### Simple server with PLAY ready connections
```go
package main

import (
	"fmt"
	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
)

func main() {
	lc := minecraft.ListenConfig{
		OnlineMode:           true,//enables server encryption
		CompressionThreshold: 0, //compresses everything!
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
```

