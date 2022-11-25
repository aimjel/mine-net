package packet

type CreateInventoryAction struct {
    Slot int16
    Present bool
    ItemId int32
    ItemCount byte
    NBT []byte

}

func (a *CreateInventoryAction) ID() int32 {
    return 0x28
}

func (a *CreateInventoryAction) Encode(w Writer) error {
    return nil
}

func (a *CreateInventoryAction) Decode(r *Reader) error {
    _ = r.Int16(&a.Slot)
    _  = r.Bool(&a.Present)
    if a.Present {
        _ = r.VarInt(&a.ItemId)
        _ = r.Uint8(&a.ItemCount)

        var nbt uint8
        r.Uint8(&nbt)
        if nbt != 0 {
            return NotImplemneted
        }
    }

    return nil
}