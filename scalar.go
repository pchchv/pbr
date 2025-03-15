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
