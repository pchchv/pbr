package pbr

import (
	"io"
	"testing"

	"github.com/pchchv/pbr/testmsg"
	"google.golang.org/protobuf/proto"
)

func TestMessage_Next(t *testing.T) {
	// read err should be false and set error
	msg := New([]byte{201, 200, 200, 200, 200, 200, 200, 200, 200, 200})
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != ErrIntOverflow {
		t.Errorf("incorrect error: %e", err)
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
		t.Errorf("incorrect error: %e", err)
	}

	// error with wire type 2, length delimited
	msg.Reset([]byte{0x10 | WireTypeLengthDelimited, 0x85, 0x04})
	msg.Next()
	msg.Skip()
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != io.ErrUnexpectedEOF {
		t.Errorf("incorrect error: %e", err)
	}

	// error with wire type 5, 32 bit
	msg.Reset([]byte{0x10 | WireType32bit, 0x85, 0x04})
	msg.Next()
	msg.Skip()
	if msg.Next() {
		t.Errorf("should be false on if error")
	}

	if err := msg.Error(); err != io.ErrUnexpectedEOF {
		t.Errorf("incorrect error: %e", err)
	}
}

func TestMessage_MessageData(t *testing.T) {
	parent := &testmsg.Parent{
		Child: &testmsg.Child{
			Number:  *proto.Int64(123),
			Numbers: []int64{1, 2, 3, -4, -5, -6, 7, 8},
			Grandchild: []*testmsg.Grandchild{
				{
					Number:  *proto.Int64(111),
					Numbers: []int64{-1, 2, -3, 4, -5, 6, -7, 8},
				},
				{
					Number:  *proto.Int64(-222),
					Numbers: []int64{1, -2, 3, -4, 5, -6, 7, -8},
				},
			},
			After: *proto.Bool(false),
		},
		After: *proto.Bool(false),
	}

	data, err := proto.Marshal(parent)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	t.Run("decode using golang/protobuf", func(t *testing.T) {
		msg := New(data)
		p := &testmsg.Parent{}

		for msg.Next() {
			switch msg.FieldNumber() {
			case 1:
				d, err := msg.MessageData()
				if err != nil {
					t.Fatalf("unable to read message: %e", err)
				}

				p.Child = &testmsg.Child{}
				if err = proto.Unmarshal(d, p.Child); err != nil {
					t.Fatalf("unable to unmarshal: %e", err)
				}
			case 32:
				v, err := msg.Bool()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.After = v
			default:
				t.Fatalf("no skips in this message: field number %d", msg.FieldNumber())
			}
		}

		if err := msg.Error(); err != nil {
			t.Fatalf("scanning error: %e", err)
		}

		compare(t, p, parent)
	})

	t.Run("invalid packed length", func(t *testing.T) {
		msg := New([]byte{200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200})
		if _, err := msg.MessageData(); err != ErrIntOverflow {
			t.Errorf("incorrect error: %e", err)
		}
	})
}

func TestMessage_Reset(t *testing.T) {
	message1 := &testmsg.Scalar{Flt: *proto.Float32(123.4567)}
	data, err := proto.Marshal(message1)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	msg := New(data)
	s := &testmsg.Scalar{}
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1:
			if v := msg.WireType(); v != WireType32bit {
				t.Errorf("incorrect wiretype: %v", v)
			}

			v, err := msg.Float()
			if err != nil {
				t.Fatalf("unable to read float: %e", err)
			}
			s.Flt = v
		default:
			msg.Skip()
		}
	}
	compare(t, s, message1)
	// Reset
	msg.Reset(nil)
	s = &testmsg.Scalar{}
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1:
			v, err := msg.Float()
			if err != nil {
				t.Fatalf("unable to read float: %e", err)
			}
			s.Flt = v
		default:
			msg.Skip()
		}
	}
	compare(t, s, message1)

	// Reset with new data
	message2 := &testmsg.Scalar{Flt: *proto.Float32(55.555)}
	data2, err := proto.Marshal(message2)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	msg.Reset(data2)
	s = &testmsg.Scalar{}
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1:
			v, err := msg.Float()
			if err != nil {
				t.Fatalf("unable to read float: %e", err)
			}
			s.Flt = v
		default:
			msg.Skip()
		}
	}

	compare(t, s, message2)
}
