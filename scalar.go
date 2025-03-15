package pbr

import (
	"encoding/binary"
	"io"
)

// base contains all methods for
// reading packable fields (numbers)
// so that they can be shared between
// the message and the iterator.
type base struct {
	Data  []byte
	Index int
}

// Fixed32 reads a fixed 4 byte value as a uint32.
// This proto type is more efficient than uint32
// if values are often greater than 2^28.
func (b *base) Fixed32() (uint32, error) {
	if len(b.Data) < b.Index+4 {
		return 0, io.ErrUnexpectedEOF
	}

	v := binary.LittleEndian.Uint32(b.Data[b.Index:])
	b.Index += 4
	return v, nil
}

// Sfixed32 reads a fixed 4 byte value signed value.
func (b *base) Sfixed32() (int32, error) {
	v, err := b.Fixed32()
	return int32(v), err
}

// Varint32 reads up to 32-bits of variable-length encoded data.
// Note that negative int32 values could still be encoded as 64-bit varints due to their leading 1s.
func (b *base) Varint32() (v uint32, err error) {
	b.Index, v, err = varint32(b.Data, b.Index)
	return
}

// Int32 reads a variable-length encoding of up to 4 bytes.
// This field type is best used if the field only has positive numbers,
// otherwise use sint32.
// Note, this field can also by read as an Int64.
func (b *base) Int32() (int32, error) {
	var v uint64
	var err error
	b.Index, v, err = varint64(b.Data, b.Index)
	return int32(v), err
}

// Uint32 reads a variable-length encoding of up to 4 bytes.
func (b *base) Uint32() (v uint32, err error) {
	b.Index, v, err = varint32(b.Data, b.Index)
	return v, err
}

// Sint32 uses variable-length encoding with
// zig-zag encoding for signed values.
// This field type more efficiently encodes
// negative numbers than regular int32s.
func (b *base) Sint32() (int32, error) {
	var v uint64
	var err error
	b.Index, v, err = varint64(b.Data, b.Index)
	return int32(unZig64(v)), err
}

func varint32(data []byte, index int) (int, uint32, error) {
	var val uint32
	shift := uint(0)
loop:
	if shift >= 32 {
		return index, 0, ErrIntOverflow
	}

	if len(data) <= index {
		return index, 0, io.ErrUnexpectedEOF
	}

	d := data[index]
	index++
	val |= uint32(d&0x7F) << shift
	if d >= 0x80 {
		shift += 7
		goto loop
	}

	return index, val, nil
}

func varint64(data []byte, index int) (int, uint64, error) {
	var val uint64
	shift := uint(0)
loop:
	if shift >= 64 {
		return 0, 0, ErrIntOverflow
	}

	if len(data) <= index {
		return 0, 0, io.ErrUnexpectedEOF
	}

	d := data[index]
	index++
	val |= uint64(d&0x7F) << shift
	if d >= 0x80 {
		shift += 7
		goto loop
	}

	return index, val, nil
}

func unZig64(v uint64) int64 {
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63))
}
