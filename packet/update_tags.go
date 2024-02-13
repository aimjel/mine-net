package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type Tag struct {
	Name    string
	Entries []int32
}

type TagType struct {
	Type string
	Tags []Tag
}

type UpdateTags struct {
	Tags []TagType
}

func (*UpdateTags) ID() int32 {
	return 0x6E
}

func (s *UpdateTags) Decode(r *encoding.Reader) error {
	var length int32
	r.VarInt(&length)
	s.Tags = make([]TagType, int(length))
	for _, t := range s.Tags {
		r.String(&t.Type)
		var length1 int32
		r.VarInt(&length1)
		t.Tags = make([]Tag, int(length1))
		for _, tag := range t.Tags {
			r.String(&tag.Name)
			var length2 int32
			r.VarInt(&length2)
			tag.Entries = make([]int32, int(length2))
			for _, e := range tag.Entries {
				r.VarInt(&e)
			}
		}
	}
	return nil
}

func (s UpdateTags) Encode(w *encoding.Writer) error {
	w.VarInt(int32(len(s.Tags)))
	for _, t := range s.Tags {
		w.String(t.Type)
		w.VarInt(int32(len(t.Tags)))
		for _, tag := range t.Tags {
			w.String(tag.Name)
			w.VarInt(int32(len(tag.Entries)))
			for _, e := range tag.Entries {
				w.VarInt(e)
			}
		}
	}
	return nil
}
