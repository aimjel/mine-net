package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"minecraft/chat"
	"os"
)

type Status struct {
	//name of the version the server is running on
	name string

	//protocol the server is compatible with
	protocol int

	//max players the server is advertised for
	max int

	//online number of players on the server
	online int32

	//samples of the server shown when the player hovers over the connection bar
	samples []string

	//description about the server
	description chat.Message

	//favicon servers icon
	favicon string
}

func NewStatus(protocol, max int, samples []string, desc string) *Status {
	s := &Status{
		name:        versionName(protocol),
		protocol:    protocol,
		max:         max,
		samples:     samples,
		description: chat.NewMessage(desc),
	}

	d, _ := s.MarshalJSON()
	if len(d) > 32767 {
		fmt.Println("server description or sample is too big")
		return nil
	}

	f, err := os.Open("server-icon.png")
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) == false {
			fmt.Printf("%v opening server-icon", err)
		}

		fmt.Printf("%v", err)
		return s
	}

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		fmt.Printf("%v decoding server-icon", err)
		return s
	}

	if cfg.Width != 64 || cfg.Height != 64 {
		fmt.Println("server icon must be 64x64")
		return s
	}

	_, _ = f.Seek(0, 0)
	m, err := png.Decode(f)
	if err != nil {
		fmt.Printf("%v decoding server icon", err)
		return s
	}

	var e png.Encoder
	e.CompressionLevel = png.DefaultCompression

	b := bytes.NewBuffer(make([]byte, 0, 1024*8))

	if err = e.Encode(b, m); err != nil {
		fmt.Printf("%v compressiong server icon", err)
		return s
	}

	if b.Len()+len(d) > 23767 {
		fmt.Printf("server icon is too big!")
		return s
	}

	s.favicon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Bytes())

	return s
}

func (s Status) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"version": map[string]interface{}{
			"name":     s.name,
			"protocol": s.protocol,
		},

		"players": map[string]interface{}{
			"max":     s.max,
			"online":  s.online,
			"samples": s.samples,
		},

		"description": s.description,
	}

	if s.favicon != "" {
		m["favicon"] = s.favicon
	}

	return json.Marshal(m)
}

func versionName(protocol int) string {
	return map[int]string{
		757: "1.18.1",
		756: "1.17.1",
		755: "1.17",
	}[protocol]
}
