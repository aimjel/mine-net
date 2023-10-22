package packet

type EntitySoundEffect struct {
	SoundID int32
  SoundName string
  HasRange bool
  Range float32
  Category int32
  EntityID int32
  Volume, Pitch float32
  Seed int64
}

func (c EntitySoundEffect) ID() int32 {
	return 0x61
}

func (c *EntitySoundEffect) Decode(r *Reader) error {
	return NotImplemented
}

func (c EntitySoundEffect) Encode(w Writer) error {
	w.VarInt(c.SoundID + 1)
  if c.SoundID + 1 = 0 {
    w.String(c.SoundName)
    if !c.HasRange {
      w.Bool(false)
    } else {
      w.Bool(true)
      w.Float32(c.Range)
    }
  }
  w.VarInt(c.Category)
  w.VarInt(c.EntityID)
  w.Float32(c.Volume)
  w.Float32(c.Pitch)
  return w.Int64(c.Seed)
}
