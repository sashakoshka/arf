package types

// Iterator is an object capable of iterating over any string-indexed map, while
// protecting its data.
type Iterator[VALUE_TYPE any] struct {
	index      int
	keys       []string
	underlying map[string] VALUE_TYPE
}

// NewIterator creates a new iterator that iterates over the specified map.
func NewIterator[VALUE_TYPE any] (
	underlying map[string] VALUE_TYPE,
) (
	iterator Iterator[VALUE_TYPE],
) {
	iterator.underlying = underlying
	iterator.keys = make([]string, len(underlying))

	index := 0
	for key, _ := range underlying {
		iterator.keys[index] = key
		index ++
	}
	
	return
}

// Key returns the current key the iterator is on.
func (iterator Iterator[VALUE_TYPE]) Key () (key string) {
	key = iterator.keys[iterator.index]
	return
}

// Value returns the current value the iterator is on.
func (iterator Iterator[VALUE_TYPE]) Value () (value VALUE_TYPE) {
	value = iterator.underlying[iterator.keys[iterator.index]]
	return
}

// Next advances the iterator by 1.
func (iterator *Iterator[VALUE_TYPE]) Next () {
	iterator.index ++
}

// End returns whether the iterator has reached the end of the map.
func (iterator Iterator[VALUE_TYPE]) End () (atEnd bool) {
	atEnd = iterator.index >= len(iterator.keys) || iterator.index < 0
	return
}
