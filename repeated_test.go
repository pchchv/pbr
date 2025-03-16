package pbr

import (
	"io"
	"log"
	"reflect"
	"testing"

	"github.com/pchchv/pbr/testmsg"
	"google.golang.org/protobuf/proto"
)

func TestDecodeRepeated_packable(t *testing.T) {
	cases := []struct {
		name    string
		skip    int
		message *testmsg.Repeated
	}{
		{
			name: "float",
			skip: 1,
			message: &testmsg.Repeated{
				Flt:   []float32{1, 1.5, 2, 2.5, -3, -3.5},
				After: *proto.Bool(true),
			},
		},
		{
			name: "double",
			skip: 2,
			message: &testmsg.Repeated{
				Dbl:   []float64{1, 1.5, 2, 2.5, -3, -3.5},
				After: *proto.Bool(true),
			},
		},
		{
			name: "int32",
			skip: 3,
			message: &testmsg.Repeated{
				I32:   []int32{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "int64",
			skip: 4,
			message: &testmsg.Repeated{
				I64:   []int64{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "uint32",
			skip: 5,
			message: &testmsg.Repeated{
				U32:   []uint32{1, 2, 3, 4, 5, 6, 2000, 3000, 4000, 5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "uint64",
			skip: 6,
			message: &testmsg.Repeated{
				U64:   []uint64{1, 2, 3, 4, 5, 6, 2000, 3000, 4000, 5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "sint32",
			skip: 7,
			message: &testmsg.Repeated{
				S32:   []int32{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "sint64",
			skip: 8,
			message: &testmsg.Repeated{
				S64:   []int64{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "fixed32",
			skip: 9,
			message: &testmsg.Repeated{
				F32:   []uint32{1, 2, 3, 4, 5, 6, 2000, 3000, 4000, 5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "fixed64",
			skip: 10,
			message: &testmsg.Repeated{
				F64:   []uint64{1, 2, 3, 4, 5, 6, 2000, 3000, 4000, 5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "sfixed32",
			skip: 11,
			message: &testmsg.Repeated{
				Sf32:  []int32{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "sfixed64",
			skip: 12,
			message: &testmsg.Repeated{
				Sf64:  []int64{1, -2, 3, -4, 5, 6, 2000, -3000, 4000, -5000},
				After: *proto.Bool(true),
			},
		},
		{
			name: "bools",
			skip: 13,
			message: &testmsg.Repeated{
				Bool:  []bool{true, true, false, false, true, false},
				After: *proto.Bool(true),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := proto.Marshal(tc.message)
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			v, _ := decodeRepeated(t, data, 0, false)
			compare(t, v, tc.message)
		})

		t.Run("skip "+tc.name, func(t *testing.T) {
			data, err := proto.Marshal(tc.message)
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			v, _ := decodeRepeated(t, data, tc.skip, false)
			compare(t, v, &testmsg.Scalar{After: *proto.Bool(true)})
		})

		t.Run("packed "+tc.name, func(t *testing.T) {
			data, err := proto.Marshal(repeatedToPacked(tc.message))
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			v, _ := decodeRepeated(t, data, 0, false)
			compare(t, v, tc.message)
		})

		t.Run("skip packed "+tc.name, func(t *testing.T) {
			data, err := proto.Marshal(repeatedToPacked(tc.message))
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			v, _ := decodeRepeated(t, data, tc.skip, false)
			compare(t, v, &testmsg.Scalar{After: *proto.Bool(true)})
		})

		t.Run("counts "+tc.name, func(t *testing.T) {
			data, err := proto.Marshal(repeatedToPacked(tc.message))
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			v, _ := decodeRepeated(t, data, 0, false)
			if len(v.Flt) != cap(v.Flt) {
				t.Errorf("incorrect counts: %v != %v", len(v.Flt), cap(v.Flt))
			}

			if len(v.Dbl) != cap(v.Dbl) {
				t.Errorf("incorrect counts: %v != %v", len(v.Dbl), cap(v.Dbl))
			}

			if len(v.I32) != cap(v.I32) {
				t.Errorf("incorrect counts: %v != %v", len(v.I32), cap(v.I32))
			}

			if len(v.I64) != cap(v.I64) {
				t.Errorf("incorrect counts: %v != %v", len(v.I64), cap(v.I64))
			}

			if len(v.U32) != cap(v.U32) {
				t.Errorf("incorrect counts: %v != %v", len(v.U32), cap(v.U32))
			}

			if len(v.U64) != cap(v.U64) {
				t.Errorf("incorrect counts: %v != %v", len(v.U64), cap(v.U64))
			}

			if len(v.S32) != cap(v.S32) {
				t.Errorf("incorrect counts: %v != %v", len(v.S32), cap(v.S32))
			}

			if len(v.F32) != cap(v.F32) {
				t.Errorf("incorrect counts: %v != %v", len(v.F32), cap(v.F32))
			}

			if len(v.F64) != cap(v.F64) {
				t.Errorf("incorrect counts: %v != %v", len(v.F64), cap(v.F64))
			}

			if len(v.Sf32) != cap(v.Sf32) {
				t.Errorf("incorrect counts: %v != %v", len(v.Sf32), cap(v.Sf32))
			}

			if len(v.Sf64) != cap(v.Sf64) {
				t.Errorf("incorrect counts: %v != %v", len(v.Sf64), cap(v.Sf64))
			}

			if len(v.Bool) != cap(v.Bool) {
				t.Errorf("incorrect counts: %v != %v", len(v.Bool), cap(v.Bool))
			}

			if len(v.Str) != cap(v.Str) {
				t.Errorf("incorrect counts: %v != %v", len(v.Str), cap(v.Str))
			}

			if len(v.Byte) != cap(v.Byte) {
				t.Errorf("incorrect counts: %v != %v", len(v.Byte), cap(v.Byte))
			}
		})

		t.Run("errors "+tc.name, func(t *testing.T) {
			data1, err := proto.Marshal(tc.message)
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			r, err := decodeRepeated(t, data1[:1], 0, true)
			if err != io.ErrUnexpectedEOF {
				log.Printf("%v", r)
				t.Errorf("incorrect error: %e", err)
			}

			data2, err := proto.Marshal(repeatedToPacked(tc.message))
			if err != nil {
				t.Fatalf("unable to marshal: %e", err)
			}

			if _, err = decodeRepeated(t, data2[:2], 0, true); err != io.ErrUnexpectedEOF {
				t.Errorf("incorrect error: %e", err)
			}
		})
	}
}

func TestDecodeRepeated_string(t *testing.T) {
	message := &testmsg.Repeated{
		Str:   []string{"hi", "", "😂"},
		After: *proto.Bool(true),
	}

	data, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	t.Run("read", func(t *testing.T) {
		msg := New(data)
		strs := []string{}
		for msg.Next() {
			if msg.FieldNumber() == 14 {
				s, err := msg.String()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				strs = append(strs, s)
			} else {
				msg.Skip()
			}
		}

		if !reflect.DeepEqual(strs, message.Str) {
			t.Logf("%v", strs)
			t.Logf("%v", message.Str)
			t.Errorf("results not equal")
		}
	})

	t.Run("skip", func(t *testing.T) {
		var found bool
		msg := New(data)
		for msg.Next() {
			if msg.FieldNumber() == 32 {
				v, err := msg.Bool()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				found = v
			} else {
				msg.Skip()
			}
		}

		if !found {
			t.Errorf("did not find after bool")
		}
	})
}

func TestDecodeRepeated_bytes(t *testing.T) {
	message := &testmsg.Repeated{
		Byte: [][]byte{
			{},
			{1, 2, 3, 4},
			{},
			{4, 3, 2, 1},
		},
		After: *proto.Bool(true),
	}

	data, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	t.Run("read", func(t *testing.T) {
		msg := New(data)
		bytes := [][]byte{}
		for msg.Next() {
			if msg.FieldNumber() == 15 {
				b, err := msg.Bytes()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				bytes = append(bytes, b)
			} else {
				msg.Skip()
			}
		}

		if !reflect.DeepEqual(bytes, message.Byte) {
			t.Logf("%v", bytes)
			t.Logf("%v", message.Byte)
			t.Errorf("results not equal")
		}
	})

	t.Run("skip", func(t *testing.T) {
		var found bool
		msg := New(data)
		for msg.Next() {
			if msg.FieldNumber() == 32 {
				v, err := msg.Bool()
				if err != nil {
					t.Fatalf("unable to read: %e", err)
				}
				found = v
			} else {
				msg.Skip()
			}
		}

		if !found {
			t.Errorf("did not find after bool")
		}
	})
}

func repeatedToPacked(r *testmsg.Repeated) *testmsg.Packed {
	return &testmsg.Packed{
		Dbl: r.Dbl, Flt: r.Flt,
		I32: r.I32, I64: r.I64,
		U32: r.U32, U64: r.U64,
		S32: r.S32, S64: r.S64,
		F32: r.F32, F64: r.F64,
		Sf32: r.Sf32, Sf64: r.Sf64,
		Bool:  r.Bool,
		Str:   r.Str,
		Byte:  r.Byte,
		After: r.After,
	}
}

func decodeRepeated(t *testing.T, data []byte, skip int, returnErr bool) (*testmsg.Repeated, error) {
	msg := New(data)
	r := &testmsg.Repeated{}
	for msg.Next() {
		if msg.FieldNumber() == skip {
			msg.Skip()
			continue
		}

		switch msg.FieldNumber() {
		case 1:
			v, err := msg.RepeatedFloat(r.Flt)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.Flt = v
		case 2:
			v, err := msg.RepeatedDouble(r.Dbl)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.Dbl = v
		case 3:
			v, err := msg.RepeatedInt32(r.I32)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.I32 = v
		case 4:
			v, err := msg.RepeatedInt64(r.I64)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.I64 = v
		case 5:
			v, err := msg.RepeatedUint32(r.U32)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.U32 = v
		case 6:
			v, err := msg.RepeatedUint64(r.U64)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.U64 = v
		case 7:
			v, err := msg.RepeatedSint32(r.S32)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.S32 = v
		case 8:
			v, err := msg.RepeatedSint64(r.S64)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.S64 = v
		case 9:
			v, err := msg.RepeatedFixed32(r.F32)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.F32 = v
		case 10:
			v, err := msg.RepeatedFixed64(r.F64)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.F64 = v
		case 11:
			v, err := msg.RepeatedSfixed32(r.Sf32)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.Sf32 = v
		case 12:
			v, err := msg.RepeatedSfixed64(r.Sf64)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.Sf64 = v
		case 13:
			v, err := msg.RepeatedBool(r.Bool)
			if err != nil {
				if returnErr {
					return nil, err
				}
				t.Fatalf("unable to read: %e", err)
			}
			r.Bool = v
		case 32:
			v, err := msg.Bool()
			if err != nil {
				t.Fatalf("unable to read after bool: %e", err)
			}
			r.After = v
		default:
			msg.Skip()
		}
	}

	if err := msg.Error(); err != nil {
		t.Fatalf("scanning error: %e", err)
	}

	return r, nil
}
