package main

import (
	_ "embed"
	"log"
)

//go:embed packet.bin
var myBytes []byte

//go:embed chunk.0.0.bin
var serverBytes []byte

func main() {
	var skip int = 18

	for i, serverByte := range serverBytes {

		if myBytes[i] != serverByte {
			if skip != 0 {
				skip--
				continue
			}
			log.Panicf("my bytes dont match server byte in index %x %v\n", i, i)
		}
	}
}
