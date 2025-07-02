package delta

type Entity interface {
	GetID() int64
	Clone() Entity
	Delta(other Entity) Delta
	ApplyDelta(d Delta) Entity
}

type Delta interface {
	ApplyTo(e Entity) Entity
}

type Player struct {
	ID   int64
	X, Y float64
	HP   int
	Name string
}

func (p *Player) GetID() int64  { return p.ID }
func (p *Player) Clone() Entity { cp := *p; return &cp }

type PlayerDelta struct {
	X    *float64
	Y    *float64
	HP   *int
	Name *string
}

func (p *Player) Delta(other Entity) Delta {
	otherP, ok := other.(*Player)
	if !ok {
		return nil // or panic
	}
	d := &PlayerDelta{}
	if p.X != otherP.X {
		d.X = &p.X
	}
	if p.Y != otherP.Y {
		d.Y = &p.Y
	}
	if p.HP != otherP.HP {
		d.HP = &p.HP
	}
	if p.Name != otherP.Name {
		d.Name = &p.Name
	}
	return d
}

func (d *PlayerDelta) ApplyTo(e Entity) Entity {
	p := e.(*Player).Clone().(*Player)
	if d.X != nil {
		p.X = *d.X
	}
	if d.Y != nil {
		p.Y = *d.Y
	}
	if d.HP != nil {
		p.HP = *d.HP
	}
	if d.Name != nil {
		p.Name = *d.Name
	}
	return p
}

func (p *Player) ApplyDelta(d Delta) Entity {
	return d.ApplyTo(p)
}
