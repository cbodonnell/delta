package example

import (
	"bytes"
	"reflect"
	"testing"
)

func TestGameState_RoundTrip(t *testing.T) {
	// Create original state with all field types
	original := &GameState{
		// Integer types
		ID:    1,
		Round: 5,
		Score: 50,
		Lives: 2,
		MaxHP: 100,

		// Floating point types
		X:     5.5,
		Y:     10.5,
		Speed: 2.5,

		// String and boolean
		PlayerName: "TestPlayer",
		IsActive:   true,

		// Slice types
		Inventory: []string{"sword", "potion"},
		Positions: []float64{1.5, 2.5, 3.5},
		PlayerIDs: []int64{10, 20, 30},
		Data:      []byte{0xAA, 0xBB, 0xCC},

		// Map types
		PlayerScores: map[string]int16{
			"alice": 100,
			"bob":   200,
		},
		ItemCounts: map[int8]int32{
			1: 5,
			2: 3,
		},
		Metadata: map[string]string{
			"level": "forest",
			"mode":  "survival",
		},
	}

	if original.GetID() != 1 {
		t.Fatalf("GetID() = %v, want 1", original.GetID())
	}

	// Clone the original (tests deep cloning)
	cloned := original.Clone().(*GameState)
	if !reflect.DeepEqual(original, cloned) {
		t.Fatalf("Clone() did not create identical copy")
	}

	// Verify it's a deep copy by modifying original
	original.Inventory[0] = "modified"
	if cloned.Inventory[0] == "modified" {
		t.Fatalf("Clone() did not create deep copy of slice")
	}
	original.Inventory[0] = "sword" // restore

	// Modify cloned state to test all change types
	cloned.Score = 75                                  // primitive change
	cloned.Lives = 1                                   // different primitive
	cloned.X = 15.5                                    // float change
	cloned.PlayerName = "ModifiedPlayer"               // string change
	cloned.IsActive = false                            // bool change
	cloned.Inventory = []string{"bow", "arrow", "map"} // slice change
	cloned.Positions = []float64{10.0, 20.0}           // different slice
	cloned.PlayerIDs = append(cloned.PlayerIDs, 40)    // slice length change
	cloned.Data = nil                                  // nil slice
	cloned.PlayerScores["alice"] = 150                 // map value change
	cloned.PlayerScores["charlie"] = 300               // map key addition
	delete(cloned.ItemCounts, 1)                       // map key deletion
	cloned.Metadata = nil                              // nil map

	// Create delta: what values are in original that differ from cloned
	delta := original.Delta(cloned)
	if delta == nil {
		t.Fatalf("Delta() returned nil")
	}

	// Apply delta to cloned state (should make it match original)
	cloned.ApplyDelta(delta)

	// Verify round-trip: cloned should now equal original
	if !reflect.DeepEqual(cloned, original) {
		t.Errorf("Round-trip failed:")
		t.Errorf("Original: %+v", original)
		t.Errorf("After delta: %+v", cloned)

		// Check specific fields for easier debugging
		if cloned.Score != original.Score {
			t.Errorf("Score mismatch: got %v, want %v", cloned.Score, original.Score)
		}
		if cloned.PlayerName != original.PlayerName {
			t.Errorf("PlayerName mismatch: got %v, want %v", cloned.PlayerName, original.PlayerName)
		}
		if !reflect.DeepEqual(cloned.Inventory, original.Inventory) {
			t.Errorf("Inventory mismatch: got %v, want %v", cloned.Inventory, original.Inventory)
		}
		if !reflect.DeepEqual(cloned.PlayerScores, original.PlayerScores) {
			t.Errorf("PlayerScores mismatch: got %v, want %v", cloned.PlayerScores, original.PlayerScores)
		}
	}

	// Test edge cases with nil delta and wrong types
	cloned.ApplyDelta(nil) // should not panic or change anything
	if !reflect.DeepEqual(cloned, original) {
		t.Errorf("ApplyDelta(nil) modified the object")
	}

	// Test delta with nil other
	if nilDelta := original.Delta(nil); nilDelta != nil {
		t.Errorf("Delta(nil) should return nil")
	}
}

func TestGameStateDelta_SerializeDeserialize(t *testing.T) {
	// Create two game states with some differences
	original := &GameState{
		ID:         1,
		Round:      5,
		Score:      100,
		X:          10.5,
		Y:          20.5,
		PlayerName: "TestPlayer",
		IsActive:   true,
		Inventory:  []string{"sword", "potion"},
		PlayerScores: map[string]int16{
			"alice": 150,
			"bob":   200,
		},
	}

	modified := &GameState{
		ID:         1,
		Round:      3,                 // different
		Score:      50,                // different
		X:          5.0,               // different
		Y:          20.5,              // same
		PlayerName: "DifferentPlayer", // different
		IsActive:   true,              // same
		Inventory:  []string{"bow"},   // different
		PlayerScores: map[string]int16{
			"alice":   100, // different value
			"charlie": 75,  // new key
		},
	}

	// Create delta from original to modified
	delta := original.Delta(modified).(*GameStateDelta)

	// Serialize the delta
	var buf bytes.Buffer
	err := delta.Serialize(&buf)
	if err != nil {
		t.Fatalf("Failed to serialize delta: %v", err)
	}

	// Deserialize into a new delta
	newDelta := &GameStateDelta{}
	err = newDelta.Deserialize(&buf)
	if err != nil {
		t.Fatalf("Failed to deserialize delta: %v", err)
	}

	// Check that the new delta matches the original
	if !reflect.DeepEqual(newDelta, delta) {
		t.Errorf("Deserialized delta does not match original:\nOriginal: %+v\nDeserialized: %+v", delta, newDelta)
	}
}
