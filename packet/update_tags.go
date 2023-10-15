package packet

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

func (*UpdateTags) Decode(*packet.Reader) error {
	return nil
}

func (s UpdateTags) Encode(w packet.Writer) error {
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
