package pbr

// Iterator allows for moving across a
// packed repeated field in a 'controlled' fashion.
type Iterator struct {
	base
	fieldNumber int
}

// Iterator will use the current field.
// Field must be a packed repeated field.
func (m *Message) Iterator(iter *Iterator) (*Iterator, error) {
	l, err := m.packedLength()
	if err != nil {
		return nil, err
	}

	if iter == nil {
		iter = &Iterator{}
	}

	iter.base = base{
		Data:  m.Data[m.Index : m.Index+l],
		Index: 0,
	}
	iter.fieldNumber = m.fieldNumber
	m.Index += l

	return iter, nil
}
