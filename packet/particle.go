package packet

type Particle struct {
    ParticleID int32
    LongDistance bool
    X, Y, Z float64
    OffsetX, OffsetY, OffsetZ float32
    ParticleData float32
    ParticleCount int32
    Data bool//varies
}

func (p *Particle) ID() int32 {
    return 0x24
}

func (e *Particle) Decode(r *Reader) error {
    return nil
}

func (p *Particle) Encode(w Writer) error {
    _ = w.Int32(p.ParticleID)
    _ = w.Bool(p.LongDistance)
    _ = w.Float64(p.X)
    _ = w.Float64(p.Y)
    _ = w.Float64(p.Z)
    _ = w.Float32(p.OffsetX)
    _ = w.Float32(p.OffsetY)
    _ = w.Float32(p.OffsetZ)
    _ = w.Float32(p.ParticleData)
    _ = w.Int32(p.ParticleCount)
    if p.Data {
        return NotImplemneted
    }
    return nil
}

