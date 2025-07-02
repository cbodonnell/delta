package delta

import (
	"io"
)

type Entity interface {
	GetID() int64
	Clone() Entity
	Delta(other Entity) Delta
	ApplyDelta(d Delta)
}

type Delta interface {
	ApplyTo(e Entity)
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
