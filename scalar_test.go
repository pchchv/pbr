package pbr

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/pchchv/pbr/testmsg"
)

func TestMessage_Varint32(t *testing.T) {
	t.Run("overflow", func(t *testing.T) {
		msg := New([]byte{230, 230, 230, 230, 230, 230})
		if _, err := msg.Varint32(); err != ErrIntOverflow {
			t.Errorf("wrong error: %e", err)
		}
	})

	t.Run("end of input", func(t *testing.T) {
		msg := New([]byte{230, 230})
		if _, err := msg.Varint32(); err != io.ErrUnexpectedEOF {
			t.Errorf("wrong error: %e", err)
		}
	})
}

func TestMessage_Varint64(t *testing.T) {
	t.Run("overflow", func(t *testing.T) {
		msg := New([]byte{230, 230, 230, 230, 230, 230, 230, 230, 230, 230})
		if _, err := msg.Varint64(); err != ErrIntOverflow {
			t.Errorf("wrong error: %e", err)
		}
	})

	t.Run("end of input", func(t *testing.T) {
		msg := New([]byte{230, 230})
		if _, err := msg.Varint64(); err != io.ErrUnexpectedEOF {
			t.Errorf("wrong error: %e", err)
		}
	})
}

func BenchmarkVarint32(b *testing.B) {
	m := New([]byte{200, 199, 198, 6, 0, 0, 0, 0})
	v, err := m.Varint32()
	if err != nil {
		b.Fatal(err)
	}

	if v != 13738952 {
		b.Fatalf("incorrect value %v != 13738952", v)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Index = 0
		if _, err := m.Varint32(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInt64(b *testing.B) {
	m := New([]byte{200, 199, 198, 6, 0, 0, 0, 0})
	v, err := m.Int64()
	if err != nil {
		b.Fatal(err)
	}

	if v != 13738952 {
		b.Fatalf("incorrect value %v != 13738952", v)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Index = 0
		if _, err := m.Int64(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBool(b *testing.B) {
	m := New([]byte{1, 0, 1, 0, 0, 0, 0, 0})
	v, err := m.Bool()
	if err != nil {
		b.Fatal(err)
	}

	if !v {
		b.Fatalf("incorrect bool")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Index = 0
		if _, err := m.Bool(); err != nil {
			b.Fatal(err)
		}
	}
}

func compare(t *testing.T, v, expected interface{}) {
	t.Helper()
	// private fields don't work with reflect.DeepEqual, so marshaling json
	vd, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	ed, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("unable to marshal: %e", err)
	}

	if !bytes.Equal(vd, ed) {
		t.Logf("%v", string(vd))
		t.Logf("%v", string(ed))
		t.Errorf("results not equal")
	}
}

func decodeScalar(t testing.TB, data []byte, skip int) *testmsg.Scalar {
	msg := New(data)
	s := &testmsg.Scalar{}
	for msg.Next() {
		if msg.FieldNumber() == skip {
			msg.Skip()
			continue
		}

		switch msg.FieldNumber() {
		case 1:
			v, err := msg.Float()
			if err != nil {
				t.Fatalf("unable to read float: %e", err)
			}
			s.Flt = v
		case 2:
			v, err := msg.Double()
			if err != nil {
				t.Fatalf("unable to read double: %e", err)
			}
			s.Dbl = v
		case 3:
			v, err := msg.Int32()
			if err != nil {
				t.Fatalf("unable to read int32: %e", err)
			}
			s.I32 = v
		case 4:
			v, err := msg.Int64()
			if err != nil {
				t.Fatalf("unable to read int64: %e", err)
			}
			s.I64 = v
		case 5:
			v, err := msg.Uint32()
			if err != nil {
				t.Fatalf("unable to read uint32: %e", err)
			}
			s.U32 = v
		case 6:
			v, err := msg.Uint64()
			if err != nil {
				t.Fatalf("unable to read uint64: %e", err)
			}
			s.U64 = v
		case 7:
			v, err := msg.Sint32()
			if err != nil {
				t.Fatalf("unable to read sint32: %e", err)
			}
			s.S32 = v
		case 8:
			v, err := msg.Sint64()
			if err != nil {
				t.Fatalf("unable to read sint64: %e", err)
			}
			s.S64 = v
		case 9:
			v, err := msg.Fixed32()
			if err != nil {
				t.Fatalf("unable to read fixed32: %e", err)
			}
			s.F32 = v
		case 10:
			v, err := msg.Fixed64()
			if err != nil {
				t.Fatalf("unable to read fixed64: %e", err)
			}
			s.F64 = v
		case 11:
			v, err := msg.Sfixed32()
			if err != nil {
				t.Fatalf("unable to read sfixed32: %e", err)
			}
			s.Sf32 = v
		case 12:
			v, err := msg.Sfixed64()
			if err != nil {
				t.Fatalf("unable to read sfixed64: %e", err)
			}
			s.Sf64 = v
		case 13:
			v, err := msg.Bool()
			if err != nil {
				t.Fatalf("unable to read bool: %e", err)
			}
			s.Bool = v
		case 14:
			v, err := msg.String()
			if err != nil {
				t.Fatalf("unable to read string: %e", err)
			}
			s.Str = v
		case 15:
			v, err := msg.Bytes()
			if err != nil {
				t.Fatalf("unable to read bytes: %e", err)
			}
			s.Byte = v
		case 32:
			v, err := msg.Bool()
			if err != nil {
				t.Fatalf("unable to read after bool: %e", err)
			}
			s.After = v
		default:
			msg.Skip()
		}
	}

	if err := msg.Error(); err != nil {
		t.Fatalf("scanning error: %e", err)
	}

	return s
}
