package packet

import "github.com/aimjel/minecraft/chat"

type ResourcePack struct {
	URL    string
	Hash   string
	Forced bool
	Prompt *chat.Message
}

func (c ResourcePack) ID() int32 {
	return 0x40
}

func (c *ResourcePack) Decode(r *Reader) error {
	//todo implement
	return NotImplemented
}

func (c ResourcePack) Encode(w *Writer) error {
	w.String(c.URL)
	w.String(c.Hash)
	w.Bool(c.Forced)
	if c.Prompt != nil {
		w.Bool(true)
		w.String(c.Prompt.String())
	} else {
		w.Bool(false)
	}
	return nil
}
