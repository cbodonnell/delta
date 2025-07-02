package delta

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type BinaryWriter struct {
	w io.Writer
}

func NewBinaryWriter(w io.Writer) *BinaryWriter {
	return &BinaryWriter{w: w}
}

func (bw *BinaryWriter) WriteByte(b byte) error {
	_, err := bw.w.Write([]byte{b})
	return err
}

func (bw *BinaryWriter) WriteBool(b bool) error {
	if b {
		return bw.WriteByte(1)
	}
	return bw.WriteByte(0)
}

func (bw *BinaryWriter) WriteInt8(v int8) error {
	return bw.WriteByte(byte(v))
}

func (bw *BinaryWriter) WriteInt16(v int16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteInt32(v int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteInt64(v int64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteUint8(v uint8) error {
	return bw.WriteByte(v)
}

func (bw *BinaryWriter) WriteUint16(v uint16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteUint32(v uint32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteUint64(v uint64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteFloat32(v float32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, math.Float32bits(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteFloat64(v float64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	_, err := bw.w.Write(buf)
	return err
}

func (bw *BinaryWriter) WriteString(s string) error {
	// Write length as varint, then string bytes
	if err := bw.WriteVarUint32(uint32(len(s))); err != nil {
		return err
	}
	_, err := bw.w.Write([]byte(s))
	return err
}

func (bw *BinaryWriter) WriteBytes(b []byte) error {
	// Write length as varint, then bytes
	if err := bw.WriteVarUint32(uint32(len(b))); err != nil {
		return err
	}
	_, err := bw.w.Write(b)
	return err
}

// Variable-length encoding for better compression
func (bw *BinaryWriter) WriteVarUint32(v uint32) error {
	for v >= 0x80 {
		if err := bw.WriteByte(byte(v) | 0x80); err != nil {
			return err
		}
		v >>= 7
	}
	return bw.WriteByte(byte(v))
}

type BinaryReader struct {
	r io.Reader
}

func NewBinaryReader(r io.Reader) *BinaryReader {
	return &BinaryReader{r: r}
}

func (br *BinaryReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(br.r, buf)
	return buf[0], err
}

func (br *BinaryReader) ReadBool() (bool, error) {
	b, err := br.ReadByte()
	return b != 0, err
}

func (br *BinaryReader) ReadInt8() (int8, error) {
	b, err := br.ReadByte()
	return int8(b), err
}

func (br *BinaryReader) ReadInt16() (int16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(buf)), nil
}

func (br *BinaryReader) ReadInt32() (int32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), nil
}

func (br *BinaryReader) ReadInt64() (int64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(buf)), nil
}

func (br *BinaryReader) ReadUint8() (uint8, error) {
	return br.ReadByte()
}

func (br *BinaryReader) ReadUint16() (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func (br *BinaryReader) ReadUint32() (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

func (br *BinaryReader) ReadUint64() (uint64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

func (br *BinaryReader) ReadFloat32() (float32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(buf)), nil
}

func (br *BinaryReader) ReadFloat64() (float64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(br.r, buf)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(buf)), nil
}

func (br *BinaryReader) ReadString() (string, error) {
	length, err := br.ReadVarUint32()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(br.r, buf)
	return string(buf), err
}

func (br *BinaryReader) ReadBytes() ([]byte, error) {
	length, err := br.ReadVarUint32()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(br.r, buf)
	return buf, err
}

func (br *BinaryReader) ReadVarUint32() (uint32, error) {
	var result uint32
	var shift uint
	for {
		b, err := br.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= uint32(b&0x7F) << shift
		if b < 0x80 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, errors.New("varint overflow")
		}
	}
	return result, nil
}
