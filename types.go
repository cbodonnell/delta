package delta

type FieldType uint8

const (
	TypeInvalid FieldType = iota
	TypeBool    FieldType = 1
	TypeInt     FieldType = 2
	TypeInt8    FieldType = 3
	TypeInt16   FieldType = 4
	TypeInt32   FieldType = 5
	TypeInt64   FieldType = 6
	TypeUint    FieldType = 7
	TypeUint8   FieldType = 8
	TypeUint16  FieldType = 9
	TypeUint32  FieldType = 10
	TypeUint64  FieldType = 11
	TypeFloat32 FieldType = 12
	TypeFloat64 FieldType = 13
	TypeString  FieldType = 14
	TypeBytes   FieldType = 15
	TypeSlice   FieldType = 16
	TypeMap     FieldType = 17
)
