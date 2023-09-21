package main

import (
	"encoding/json"
	"errors"
	"github.com/aimjel/minecraft"
	"github.com/aimjel/minecraft/packet"
	"io"
	"os"
	"strconv"
)

func main() {
	d := minecraft.Dialer{Username: "Aimjel"}
	c, err := d.Dial("localhost:25565")
	if err != nil {
		panic(err)
	}

	m := make(map[string]statistics)
	bundles := make([][]int32, 0)
	bundle := make([]int32, 0)
	var inBundle bool
	for {
		pk, err := c.ReadPacket()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			panic(err)
		}

		if uk, ok := pk.(packet.Unknown); ok {
			id := strconv.FormatInt(int64(uk.Id), 16)
			stat := m[id]

			stat.Count++
			stat.TotalBytes += len(uk.Payload)

			if stat.MinSize == 0 {
				stat.MinSize = len(uk.Payload)
			}

			if len(uk.Payload) < stat.MinSize {
				stat.MinSize = len(uk.Payload)
			}

			if len(uk.Payload) > stat.MaxSize {
				stat.MaxSize = len(uk.Payload)
			}

			stat.entries = append(stat.entries, len(uk.Payload))

			m[id] = stat

			if uk.Id == 0x00 {
				inBundle = !inBundle

				if inBundle == false {
					bundles = append(bundles, bundle)
					bundle = bundle[:0]
				}
				continue
			}

			if inBundle {
				bundle = append(bundle, uk.Id)
			}
		}
	}

	stats, _ := json.MarshalIndent(m, "", "    ")
	os.WriteFile("stats.json", stats, 0666)
}

type statistics struct {
	Count int

	TotalBytes int

	MinSize int

	MaxSize int

	entries []int
}
