package packet


type SetSlot struct {
    WindowID byte
    StateID int32
    Slot int16
    Present bool
    ItemId int32
    ItemCount byte
    NBT []byte
}

func (s *SetSlot) ID() int32 {
    return 0x16
}

func (s *SetSlot) Encode(w Writer) error {
    _ = w.Uint8(s.WindowID)
    _ = w.VarInt(s.StateID)
    _ = w.Int16(s.Slot)
    _  = w.Bool(s.Present)
    if s.Present {
        _ = w.VarInt(s.ItemId)
        _ = w.Uint8(s.ItemCount)
        return w.Nbt(s.NBT)
    }

    return nil
}