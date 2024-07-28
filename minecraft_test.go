package minecraft

import (
	"fmt"
	"testing"
)

func TestListenConfig_Listen(t *testing.T) {
	cfg := ListenConfig{
		Status:               nil,
		OnlineMode:           false,
		CompressionThreshold: 256,
	}

	ln, err := cfg.Listen("localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		fmt.Printf("welcome %v", conn.RemoteAddr().String())
	}
}
