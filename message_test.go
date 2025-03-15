package pbr

import (
	"io"
	"testing"
)

func TestMessage_Next(t *testing.T) {
	// read err should be false and set error
	msg := New([]byte{201, 200, 200, 200, 200, 200, 200, 200, 200, 200})
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != ErrIntOverflow {
		t.Errorf("incorrect error: %v", err)
	}
}

func TestMessage_Skip(t *testing.T) {
	// error with wire type 1, 64 bit
	msg := New([]byte{0x10 | WireType64bit, 0x05})
	msg.Next()
	msg.Skip()
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != io.ErrUnexpectedEOF {
		t.Errorf("incorrect error: %v", err)
	}

	// error with wire type 2, length delimited
	msg.Reset([]byte{0x10 | WireTypeLengthDelimited, 0x85, 0x04})
	msg.Next()
	msg.Skip()
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != io.ErrUnexpectedEOF {
		t.Errorf("incorrect error: %v", err)
	}

	// error with wire type 5, 32 bit
	msg.Reset([]byte{0x10 | WireType32bit, 0x85, 0x04})
	msg.Next()
	msg.Skip()
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != io.ErrUnexpectedEOF {
		t.Errorf("incorrect error: %v", err)
	}
}
