package minecraft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"minecraft/packet"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	c, err := Dial("mcobject.tk:25560")
	if err != nil {
		t.Fatal(err)
	}

	if err = c.WritePacket(&packet.Handshake{
		ProtocolVersion: 756,
		NextState:       2,
	}); err != nil {
		t.Fatal(err)
	}

	if err = c.WritePacket(&packet.LoginStart{Name: "Aimjel"}); err != nil {
		t.Fatal(err)
	}

	c.pool = map[int32]func() packet.Packet{
		0x12: func() packet.Packet { return &packet.DeclareCommands{} },
		0x36: func() packet.Packet { return &packet.PlayerInfo{} },
	}
	for {
		p, err := c.ReadPacket()
		if err != nil {
			t.Fatal(err)
		}

		if p != nil {
			if p.ID() == 0x12 {
				b, _ := json.MarshalIndent(p, "", "	")
				os.WriteFile("declare_commands.json", b, 0666)
				fmt.Printf("%v\n", p)
			}
		}
	}
}

func BenchmarkConn_ReadPacket(b *testing.B) {
	c := newConn(nil)
	c.dec.rd = bytes.NewBuffer(make([]byte, 20))
	c.enc.wr = bytes.NewBuffer(make([]byte, 0, 20))
	c.pool = handshakePool

	if err := c.WritePacket(&packet.Handshake{
		ProtocolVersion: 760,
		ServerAddress:   "localhost",
		ServerPort:      25565,
		NextState:       1,
	}); err != nil {
		b.Fatal(err)
	}

	c.dec.w = copy(c.dec.buf, c.enc.wr.(*bytes.Buffer).Bytes())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.ReadPacket()
		if err != nil {
			b.Fatal(err)
		}

		c.dec.r = 0
	}

	b.ReportAllocs()
}

func BenchmarkConn_DecodePacket(b *testing.B) {
	c := newConn(nil)
	c.dec.rd = bytes.NewBuffer(make([]byte, 20))
	c.enc.wr = bytes.NewBuffer(make([]byte, 0, 20))

	if err := c.WritePacket(&packet.Handshake{
		ProtocolVersion: 760,
		ServerAddress:   "localhost",
		ServerPort:      25565,
		NextState:       1,
	}); err != nil {
		b.Fatal(err)
	}

	c.dec.w = copy(c.dec.buf, c.enc.wr.(*bytes.Buffer).Bytes())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var pk packet.Handshake
		if err := c.DecodePacket(&pk); err != nil {
			b.Fatal(err)
		}

		c.dec.r = 0
	}

	b.ReportAllocs()
}

func TestSpamDial(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			err := createBot("bottt" + strconv.FormatInt(int64(i), 10))
			if err != nil {
				fmt.Printf("%v %v\n", "bottt"+strconv.FormatInt(int64(i), 10), err)
				wg.Done()
			}
		}(i, &wg)
	}

	wg.Wait()
}

func createBot(name string) error {
	c, err := Dial("172.93.105.106:25568")
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.WritePacket(&packet.Handshake{
		ProtocolVersion: 47, //1.8.9
		NextState:       2,
	}); err != nil {
		return err
	}

	if err = c.WritePacket(&packet.LoginStart{Name: name}); err != nil {
		return err
	}

	c.pool = map[int32]func() packet.Packet{
		0x08: func() packet.Packet { return &packet.PlayerPositionLook{} },
	}

	var pos packet.PlayerPositionLook

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			c.Close()
			return err
		}

		if p, ok := pk.(*packet.PlayerPositionLook); ok {
			pos = *p

			//c.WritePacket(&packet.TeleportConfirm{TeleportID: p.TeleportID})
			break
		}
	}

	pk := packet.PlayerPositionRotation{
		X:     pos.X,
		FeetY: pos.Y,
		Z:     pos.Z,
	}
	pk.FeetY += 20 //test for my server
	//var n int

	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		//fmt.Printf("%v sending pk\n", name)
		pk.X += s.Float64() - s.Float64()
		pk.Z -= s.Float64() - s.Float64()
		if err = c.WritePacket(&pk); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
}
