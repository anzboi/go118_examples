package stream

import "go1_18/generics/iterator"

// stream[T, S] implements Stream[S]
//
// The base underlying values come from an iterator[T]
// The overlay allows us to intercept and apply transformations we want before processing items in the stream. The stream type is Stream[S] for this reason
type stream[T, S any] struct {
	base    iterator.Iterator[T]
	overlay func(T) S // overlay allows us to intercept values coming out of the stream
}

func NewStream[T any](elems []T) Stream[T] {
	return &stream[T, T]{
		base:    iterator.NewIterator(elems),
		overlay: func(t T) T { return t },
	}
}

func (s *stream[T, S]) ForEach(fn func(S)) {
	for s.base.Next() {
		fn(s.overlay(s.base.Value()))
	}
}

func (s *stream[T, S]) Filter(fn func(S) bool) Stream[S] {
	filterFn := func(t T) bool {
		return fn(s.overlay(t))
	}
	return &stream[T, S]{
		base:    &filterIter[T]{Iterator: s.base, filterFn: filterFn},
		overlay: s.overlay,
	}
}

// Count counts the number of elements in the stream
// This is destructive.
func (s *stream[T, S]) Count() int {
	count := 0
	s.ForEach(func(s S) {
		count++
	})
	return count
}

// AsList converts the stream into a list
// Destructive
func (s *stream[T, S]) AsList() []S {
	ret := make([]S, 0)
	s.ForEach(func(item S) {
		ret = append(ret, item)
	})
	return ret
}

// First returns the first element in the stream
// Destructive in the sense that the stream no longer
// has access to the first element
func (s *stream[T, S]) First() (S, bool) {
	if !s.base.Next() {
		var zero S
		return zero, false
	}
	return s.overlay(s.base.Value()), true
}

// This must be defined as a package method. We cannot define Stream[T].Map[S], no generic methods.
//
// This introduces a new problem though. We must implement a new stream on top of the given one without access to the stream implementation.
// The mapFn becomes the overlay function in the new stream (and is the primary reason we have the overlay function in the first place)
func Map[T, S any](s Stream[T], mapFn func(T) S) Stream[S] {
	// TODO
	return nil
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
