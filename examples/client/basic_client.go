//go:build ignore

package main

import (
	"fmt"
	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
)

func main() {
	d := minecraft.Dialer{Username: "Aimjel"}
	c, err := d.Dial("localhost:25565")
	if err != nil {
		panic(err)
	}

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			panic(err)
		}

		if uk, ok := pk.(packet.Unknown); ok {
			if uk.Id == 0x3a {
				fmt.Println(uk.Payload)
			}
			fmt.Println(uk.Id, len(uk.Payload))
		}
	}
}
