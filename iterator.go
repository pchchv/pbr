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

// Count returns the total number of values in the given repeating field.
// The answer depends on the type/encoding or the field:
// double, float, fixed, sfixed are WireType32bit or WireType64bit,
// all other types (int, uint, sint) are WireTypeVarint.
// Any other value will cause the function to panic.
func (i *Iterator) Count(wireType int) (count int) {
	switch wireType {
	case WireType32bit:
		return len(i.base.Data) / 4
	case WireType64bit:
		return len(i.base.Data) / 8
	case WireTypeVarint:
		for _, b := range i.Data {
			if b < 128 {
				count++
			}
		}
		return
	default:
		panic("invalid wire type for a packed repeated field")
	}
}

// Skip will move the interator forward 'count' value without actually reading it.
// For a new iterator,
// 'count' will move the pointer so that the next value call will be the 'counth' value.
// The correct wireType must be specified:
// double, float, fixed, sfixed are WireType32bit or WireType64bit,
// all other types (int, uint, sint) are WireTypeVarint.
// Any other value will cause the function to panic.
func (i *Iterator) Skip(wireType int, count int) {
	switch wireType {
	case WireTypeVarint:
		for j := 0; j < count; j++ {
			for i.Data[i.Index] >= 128 {
				i.Index++
			}
			i.Index++
		}
	case WireType32bit:
		i.Index += 4 * count
	case WireType64bit:
		i.Index += 8 * count
	default:
		panic("invalid wire type for a packed repeated field")
	}
}
