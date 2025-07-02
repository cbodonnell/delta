package example

// syncgen:entity
type GameState struct {
	// Integer types
	ID    int64
	Round int32
	Score int
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
	PlayerScores map[string]int
	ItemCounts   map[int]int
	Metadata     map[string]string
}
