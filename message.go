package pbr

import (
	"errors"
	"io"
)

var (
	// ErrIntOverflow is returned when scanning a varint-encoded integer,
	// the value is found to be too long for the integer type.
	ErrIntOverflow = errors.New("protoscan: integer overflow")
	// ErrInvalidLength is returned when the length is not valid,
	// usually as a result of an invalid type scan.
	ErrInvalidLength = errors.New("protoscan: invalid length")
)

// Message is a container for a protobuf message type ready to be scanned.
type Message struct {
	base
	err         error
	fieldNumber int
	wireType    int
}

// New creates a new Message scanner for the given encoded protobuf data.
func New(data []byte) *Message {
	return &Message{
		base: base{
			Data:  data,
			Index: 0,
		},
	}
}

// Reset will set the index to 0 so the message can be read again.
// Optionally pass in new data to reuse the Message object.
func (m *Message) Reset(newData []byte) {
	if newData != nil {
		m.Data = newData
	}

	m.err = nil
	m.Index = 0
	m.fieldNumber = 0
	m.wireType = 0
}

// Next will move the scanner to the next value.
// Should be used in a for loop.
func (m *Message) Next() bool {
	if m.err == nil && m.Index < len(m.Data) {
		if val, err := m.Varint64(); err != nil {
			m.err = err
			return false
		} else {
			m.fieldNumber = int(val >> 3)
			m.wireType = int(val & 0x7)
			return true
		}
	}

	return false
}

func (m *Message) packedLength() (l int, err error) {
	var l64 uint64
	m.Index, l64, err = varint64(m.Data, m.Index)
	if err != nil {
		return
	}

	l = int(l64)
	if l < 0 {
		return 0, ErrInvalidLength
	}

	postIndex := m.Index + l
	if postIndex < 0 {
		return 0, ErrInvalidLength
	}

	if len(m.Data) < postIndex {
		return 0, io.ErrUnexpectedEOF
	}

	return
}
