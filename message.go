package pbr

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
