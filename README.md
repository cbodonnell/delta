# Delta - Game State Synchronization

Efficient delta compression for multiplayer games in Go. Only transmit what changed, not the entire state.

## Installation

```bash
go get github.com/cbodonnell/delta
go install github.com/cbodonnell/delta/cmd/deltagen@latest
```

## Quick Start

### 1. Define Your Game State

```go
// delta:entity
type GameState struct {
    ID         int64   // Required
    PlayerX    float64
    PlayerY    float64
    Score      int32
    PlayerName string
    Inventory  []string
}
```

### 2. Generate Delta Code

```bash
deltagen -input .
```

### 3. Use Deltas

```go
// Create states
state1 := &GameState{ID: 1, PlayerX: 10, Score: 100}
state2 := state1.Clone().(*GameState)
state2.PlayerX = 15  // Player moved
state2.Score = 120   // Score changed

// Create and serialize delta
delta := state1.Delta(state2)
var buf bytes.Buffer
delta.Serialize(&buf)

// Apply delta elsewhere
newDelta := &GameStateDelta{}
newDelta.Deserialize(&buf)
state1.ApplyDelta(newDelta)  // state1 now equals state2
```

## Supported Types

- **Primitives**: `bool`, `int8`-`int64`, `uint8`-`uint64`, `float32`, `float64`, and `string`
- **Collections**: `[]T`, `map[K]V`, and `[]byte`, where `K` and `V` are supported primitive types

## Network Usage

```go
// Server: send only changes
delta := client.lastState.Delta(newState)
sendToClient(delta)
client.lastState = newState.Clone()

// Client: apply changes
delta := receiveDelta()
gameState.ApplyDelta(delta)
```

## Requirements

- Structs must have `// delta:entity` comment
- Must include `ID int64` field
- Only exported fields are processed