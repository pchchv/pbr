package pbr

// Iterator allows for moving across a
// packed repeated field in a 'controlled' fashion.
type Iterator struct {
	base
	fieldNumber int
}
