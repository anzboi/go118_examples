package main

type set[T comparable] map[T]struct{}

func (s set[T]) Has(elem T) bool {
	_, ok := s[elem]
	return ok
}

func (s set[T]) Add(elem T) {
	s[elem] = struct{}{}
}

type IntSet struct {
	set[int]
}
