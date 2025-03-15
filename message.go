package pbr

// base contains all methods for
// reading packable fields (numbers)
// so that they can be shared between
// the message and the iterator.
type base struct {
	Data  []byte
	Index int
}

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
