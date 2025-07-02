package example

// delta:entity
type GameState struct {
	// Integer types
	ID    int64
	Round int16
	Score int32
	Lives int8
	MaxHP uint16

	// Floating point types
	X, Y  float64
	Speed float32

	// String
	PlayerName string

	// Boolean
	IsActive bool

	// Slice types
	Inventory []string
	Positions []float64
	PlayerIDs []int64
	Data      []byte

	// Map types
	PlayerScores map[string]int16
	ItemCounts   map[int8]int32
	Metadata     map[string]string
}
