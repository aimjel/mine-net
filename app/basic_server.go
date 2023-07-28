package main

import (
	"fmt"
	"github.com/aimjel/minecraft"
)

func main() {

	lc := minecraft.ListenConfig{
		OnlineMode:           true,
		CompressionThreshold: 0, //compresses everything!
		Status:               minecraft.NewStatus(756, 10, "someone had todo it"),
	}

	l, _ := lc.Listen(":25565")

	for {
		_, err := l.Accept()
		if err != nil {
			fmt.Printf("main thread|server: %v", err)
			return
		}

		//c is a connection ready to join the world
	}
}
