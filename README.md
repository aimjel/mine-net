# Minecraft Protocol Library for Go

Efficient Minecraft protocol integration in Go. 
Send, receive, and manipulate packets with its built-in connection class. 
TCP listener manages status and login states, complete with full encryption and compression support.

## Installation

```sh
go get github.com/aimjel/minenet
```

### Simple server with PLAY ready connections
```go
package main

import (
	"fmt"
	"github.com/aimjel/minenet"
	"github.com/aimjel/minenet/packet"
)

func main() {
	lc := minenet.ListenConfig{
		OnlineMode:           true,//enables server encryption
		CompressionThreshold: 0, //compresses everything!
		Status:               minenet.NewStatus(minecraft.Version{Protocol: 763}, 10, "someone had todo it"),
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

