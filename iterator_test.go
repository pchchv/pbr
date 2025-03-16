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

func decodeIterator(t *testing.T, data []byte, skip int) *testmsg.Packed {
	msg := New(data)
	p := &testmsg.Packed{}
	for msg.Next() {
		if msg.FieldNumber() == skip {
			msg.Skip()
			continue
		}

		switch msg.FieldNumber() {
		case 1:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.Flt = make([]float32, 0, iter.Count(WireType32bit))
			for iter.HasNext() {
				v, err := iter.Float()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.Flt = append(p.Flt, v)
			}
		case 2:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.Dbl = make([]float64, 0, iter.Count(WireType64bit))
			for iter.HasNext() {
				v, err := iter.Double()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.Dbl = append(p.Dbl, v)
			}
		case 3:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.I32 = make([]int32, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Int32()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.I32 = append(p.I32, v)
			}
		case 4:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.I64 = make([]int64, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Int64()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.I64 = append(p.I64, v)
			}
		case 5:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.U32 = make([]uint32, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Uint32()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.U32 = append(p.U32, v)
			}
		case 6:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.U64 = make([]uint64, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Uint64()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.U64 = append(p.U64, v)
			}
		case 7:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.S32 = make([]int32, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Sint32()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.S32 = append(p.S32, v)
			}
		case 8:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.S64 = make([]int64, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Sint64()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.S64 = append(p.S64, v)
			}
		case 9:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.F32 = make([]uint32, 0, iter.Count(WireType32bit))
			for iter.HasNext() {
				v, err := iter.Fixed32()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.F32 = append(p.F32, v)
			}
		case 10:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.F64 = make([]uint64, 0, iter.Count(WireType64bit))
			for iter.HasNext() {
				v, err := iter.Fixed64()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.F64 = append(p.F64, v)
			}
		case 11:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.Sf32 = make([]int32, 0, iter.Count(WireType32bit))
			for iter.HasNext() {
				v, err := iter.Sfixed32()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.Sf32 = append(p.Sf32, v)
			}
		case 12:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.Sf64 = make([]int64, 0, iter.Count(WireType64bit))
			for iter.HasNext() {
				v, err := iter.Sfixed64()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.Sf64 = append(p.Sf64, v)
			}
		case 13:
			iter, err := msg.Iterator(nil)
			if err != nil {
				t.Fatalf("unable to create iterator: %e", err)
			}

			p.Bool = make([]bool, 0, iter.Count(WireTypeVarint))
			for iter.HasNext() {
				v, err := iter.Bool()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				p.Bool = append(p.Bool, v)
			}
		case 32:
			v, err := msg.Bool()
			if err != nil {
				t.Fatalf("unable to read after bool: %e", err)
			}
			p.After = v
		default:
			msg.Skip()
		}
	}

	if err := msg.Error(); err != nil {
		t.Fatalf("scanning error: %e", err)
	}

	return p
}
