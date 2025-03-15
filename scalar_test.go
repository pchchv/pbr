package pbr

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
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
