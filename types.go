package delta

type FieldType uint8

// Comprehensive type system for delta compression
const (
	TypeInvalid FieldType = iota

	// Boolean types
	TypeBool

	// Integer types
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64

	// Floating point types
	TypeFloat32
	TypeFloat64

	// String and byte types
	TypeString

	// Collection types
	TypeSlice
	TypeMap
)

// TypeInfo contains metadata about a field type
type TypeInfo struct {
	Type FieldType
	Name string
}

// GetTypeInfo analyzes a type string and returns TypeInfo
func GetTypeInfo(typeStr string) TypeInfo {
	info := TypeInfo{Name: typeStr}

	switch typeStr {
	case "bool":
		info.Type = TypeBool
	case "int8":
		info.Type = TypeInt8
	case "int16":
		info.Type = TypeInt16
	case "int", "int32":
		info.Type = TypeInt32
	case "int64":
		info.Type = TypeInt64
	case "uint8", "byte":
		info.Type = TypeUint8
	case "uint16":
		info.Type = TypeUint16
	case "uint", "uint32":
		info.Type = TypeUint32
	case "uint64":
		info.Type = TypeUint64
	case "float32":
		info.Type = TypeFloat32
	case "float64":
		info.Type = TypeFloat64
	case "string":
		info.Type = TypeString
	default:
		// Analyze complex types
		switch {
		case len(typeStr) >= 2 && typeStr[:2] == "[]":
			info.Type = TypeSlice
		case len(typeStr) >= 4 && typeStr[:4] == "map[":
			info.Type = TypeMap
		}
	}

	return info
}

// Helper functions for slice comparison
func SlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper functions for map comparison
func MapsEqual[K, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || v != bv {
			return false
		}
	}
	return true
}
