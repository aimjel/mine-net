package packet

import "github.com/aimjel/minecraft/chat"

type ResourcePack struct {
	URL string
  Hash string
  Forced bool
  Prompt string
}

func (c ResourcePack) ID() int32 {
	return 0x40
}

func (c *ResourcePack) Decode(r *Reader) error {
  //todo implement
	return nil
}

func (c ResourcePack) Encode(w Writer) error {
  msg := chat.NewMessage(c.Prompt)
	w.String(c.URL)
  w.String(c.Hash)
  w.Bool(c.Forced)
  if c.Prompt != "" {
    w.Bool(true)
    w.String(msg.String())
  } else {
    w.Bool(false)
  }
	return nil
}
