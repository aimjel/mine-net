package protocol

import (
	"bytes"
	"github.com/aimjel/minecraft/packet"
	"testing"
)

var hs = &packet.Handshake{
	ProtocolVersion: 758,
	ServerAddress:   "localhost",
	ServerPort:      25565,
	NextState:       1,
}

var jg = &packet.JoinGame{}

func TestEncoder_Encode(t *testing.T) {
	enc := NewEncoder()
	//enc.EnableCompression(1)

	for i := 0; i < 1; i++ {
		if err := enc.Encode(jg); err != nil {
			t.Fatal(err)
		}

	}

	//t.Log(enc.Flush())
}

func TestDecoder_DecodePacket(t *testing.T) {
	t.Run("TestEncoder_Encode", TestEncoder_Encode)

	enc := NewEncoder()
	enc.EnableCompression(0)

	for i := 0; i < 1; i++ {
		if err := enc.Encode(hs); err != nil {
			t.Fatal(err)
		}
	}

	dec := NewDecoder(bytes.NewBuffer(enc.Flush()))
	dec.EnableDecompression()

	for i := 0; i < 1; i++ {
		pk, err := dec.DecodePacket()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v\n", pk)
	}
}
