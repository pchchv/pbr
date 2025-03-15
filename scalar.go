package pbr

// base contains all methods for
// reading packable fields (numbers)
// so that they can be shared between
// the message and the iterator.
type base struct {
	Data  []byte
	Index int
}
