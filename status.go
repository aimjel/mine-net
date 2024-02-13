package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"image/png"
	"math"
	"os"

	"github.com/aimjel/minecraft/chat"
)

type Version struct {
	Protocol int
	Text     string
}

type Status struct {
	enc *json.Encoder

	buf *bytes.Buffer

	s *status
}

func NewStatus(version Version, max int, desc string, enforcesSecureChat, previewsChat bool) *Status {
	var s status
	if version.Text == "" {
		version.Text = versionName(version.Protocol)
	}
	s.Version.Name, s.Version.Protocol = version.Text, version.Protocol
	s.Players.Max, s.Description = max, chat.NewMessage(desc)
	s.EnforcesSecureChat, s.PreviewsChat = enforcesSecureChat, previewsChat

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	_ = enc.Encode(s)

	st := &Status{enc: enc, buf: &buf, s: &s}

	if err := st.loadIcon(); err != nil {
		return st
	}

	return st
}

func (s *Status) loadIcon() error {
	f, err := os.Open("server-icon.png")
	if err != nil {
		return err
	}
	defer f.Close()

	m, err := png.Decode(f)
	if err != nil {
		return err
	}

	if b := m.Bounds(); b.Dy() != 64 || b.Dx() != 64 {
		return fmt.Errorf("server-icon must be a 64x64")
	}

	var e png.Encoder
	e.CompressionLevel = png.DefaultCompression

	statusLen := s.buf.Len()
	s.buf.Reset()

	if err = e.Encode(s.buf, m); err != nil {
		fmt.Printf("%v compressiong server icon", err)
	}

	prependedInfo := "\"favicon\": \"data:image/png;base64,\","
	iconLen := len(prependedInfo) + base64.StdEncoding.EncodedLen(s.buf.Len())
	if statusLen+iconLen > math.MaxInt16 {
		return fmt.Errorf("server-icon file is too big")
	}

	s.s.Favicon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(s.buf.Bytes())
	s.buf.Reset()

	return s.enc.Encode(s.s)
}

func (s *Status) json() []byte {
	return s.buf.Bytes()
}

func (s *Status) handleStatus(c *Conn) error {
	var rq packet.Request
	if err := c.DecodePacket(&rq); err != nil {
		return err
	}

	if err := c.SendPacket(&packet.Response{JSON: s.json()}); err != nil {
		return fmt.Errorf("%v writing response packet", err)
	}

	var pg packet.Ping
	if err := c.DecodePacket(&pg); err != nil {
		return fmt.Errorf("%v decoding ping packet", err)
	}

	return c.SendPacket(&packet.Pong{Payload: pg.Payload})
}

func versionName(protocol int) string {
	return map[int]string{
		763: "1.20/1.20.1",
		762: "1.19.4",
		761: "1.19.3",
		760: "1.19.1/1.19.2",
		759: "1.19",
		758: "1.18.2",
		757: "1.18.1",
		756: "1.17.1",
		755: "1.17",
	}[protocol]
}

// status represents the json response in struct form for better performance
type status struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description        chat.Message `json:"description"`
	Favicon            string       `json:"favicon,omitempty"`
	EnforcesSecureChat bool         `json:"enforcesSecureChat"`
	PreviewsChat       bool         `json:"previewsChat"`
}
