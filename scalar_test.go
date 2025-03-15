package pbr

import (
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
