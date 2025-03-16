package pbr

import (
	"io"
	"testing"

	"github.com/pchchv/pbr/testmsg"
	"google.golang.org/protobuf/proto"
)

func TestIterator_errors(t *testing.T) {
	message := &testmsg.Packed{
		I64: make([]int64, 4000),
	}
	data, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	msg := New(data[:2])
	if !msg.Next() {
		t.Fatalf("next is false?")
	}

	_, err = msg.Iterator(nil)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("incorrect error: %e", err)
	}
}

func TestIterator_Skip(t *testing.T) {
	message := &testmsg.Packed{
		Flt: make([]float32, 10),
		Dbl: make([]float64, 10),
		I32: make([]int32, 10),
		I64: make([]int64, 10),
	}

	for i := 0; i < 10; i++ {
		message.Flt[i] = float32(10 * i)
		message.Dbl[i] = float64(15 * i)
		message.I32[i] = int32(10 + 100*i)
	}

	message.I64[0] = int64(1 << 7)
	message.I64[2] = int64(1 << 15)
	message.I64[4] = int64(1 << 23)
	message.I64[6] = int64(1 << 23)

	data, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	msg := New(data)
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1: // Float
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to make iterator: %e", err)
			}

			iter.Skip(WireType32bit, 0)
			if v, _ := iter.Float(); v != 0 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireType32bit, 1)
			if v, _ := iter.Float(); v != 20 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireType32bit, 2)
			if v, _ := iter.Float(); v != 50 {
				t.Errorf("incorrect value: %v", v)
			}
		case 2:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to make iterator: %e", err)
			}

			iter.Skip(WireType64bit, 1)
			if v, _ := iter.Double(); v != 15 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireType64bit, 1)
			if v, _ := iter.Double(); v != 45 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireType64bit, 2)
			if v, _ := iter.Double(); v != 90 {
				t.Errorf("incorrect value: %v", v)
			}
		case 3:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to make iterator: %e", err)
			}

			iter.Skip(WireTypeVarint, 1)
			if v, _ := iter.Int32(); v != 110 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireTypeVarint, 1)
			if v, _ := iter.Int32(); v != 310 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireTypeVarint, 2)
			if v, _ := iter.Int32(); v != 610 {
				t.Errorf("incorrect value: %v", v)
			}
		case 4:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to make iterator: %e", err)
			}

			iter.Skip(WireTypeVarint, 0)
			if v, _ := iter.Int64(); v != 128 {
				t.Errorf("incorrect value: %v", v)
			}

			iter.Skip(WireTypeVarint, 2)
			if v, _ := iter.Int64(); v != 0 {
				t.Errorf("incorrect value: %v", v)
			}

			if v, _ := iter.Int64(); v != 0x0800000 {
				t.Errorf("incorrect value: %x", v)
			}
		default:
			msg.Skip()
		}
	}

	if err := msg.Error(); err != nil {
		t.Fatalf("read error: %e", err)
	}
}
