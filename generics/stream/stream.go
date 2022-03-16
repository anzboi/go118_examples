package stream

import "github.com/anzboi/go118-examples/generics/iterator"

type Stream[T any] interface {
	ForEach(func(T))
	Filter(func(T) bool) Stream[T]
	Count() int
	AsList() []T
	First() (T, bool)

	isStream() // No one else is allowed to implement
}

// stream implements the Stream interface
//
// *stream[T] implements Stream[T]
type stream[T any] struct {
	base iterator.Iterator[T]
}

func (s *stream[T]) isStream() {}

// Returns a new stream from the given iterator
func NewStream[T any](iter iterator.Iterator[T]) Stream[T] {
	return &stream[T]{
		base: iter,
	}
}

func NewStreamOf[T any](elems ...T) Stream[T] {
	return &stream[T]{
		base: iterator.NewIterator(elems),
	}
}

func (s *stream[T]) ForEach(fn func(T)) {
	for s.base.Next() {
		fn(s.base.Value())
	}
}

func (s *stream[T]) Filter(fn func(T) bool) Stream[T] {
	filterFn := func(t T) bool {
		return fn(t)
	}
	return &stream[T]{
		base: &filterIter[T]{Iterator: s.base, filterFn: filterFn},
	}
}

// Count counts the number of elements in the stream
// Destructive
func (s *stream[T]) Count() int {
	count := 0
	s.ForEach(func(T) {
		count++
	})
	return count
}

// AsList converts the stream into a list
// Destructive
func (s *stream[T]) AsList() []T {
	ret := make([]T, 0)
	s.ForEach(func(item T) {
		ret = append(ret, item)
	})
	return ret
}

// First returns the first element in the stream
// Destructive in the sense that the stream no longer
// has access to the first element
func (s *stream[T]) First() (T, bool) {
	if !s.base.Next() {
		var zero T
		return zero, false
	}
	return s.base.Value(), true
}

// This must be defined as a package method. We cannot define Stream[T].Map[S], no generic methods.
//
// This introduces a new problem though. We must implement a new stream on top of the given one without access to the stream implementation.
// The mapFn becomes the overlay function in the new stream (and is the primary reason we have the overlay function in the first place)
func Map[T, S any](s Stream[T], mapFn func(T) S) Stream[S] {
	// We can assume the only implementation of Stream is *stream
	impl, ok := s.(*stream[T])
	if !ok {
		panic("Unrecognised stream implementation")
	}
	return &stream[S]{
		base: &mapIter[T, S]{
			Iterator: impl.base,
			mapFn:    mapFn,
		},
	}
}

type filterIter[T any] struct {
	iterator.Iterator[T]
	filterFn func(T) bool
}

func (f *filterIter[T]) Next() bool {
	for f.Iterator.Next() {
		if f.filterFn(f.Iterator.Value()) {
			return true
		}
	}
	return false
}

type mapIter[T, S any] struct {
	iterator.Iterator[T]
	mapFn func(T) S
}

func (m *mapIter[T, S]) Value() S {
	return m.mapFn(m.Iterator.Value())
}
