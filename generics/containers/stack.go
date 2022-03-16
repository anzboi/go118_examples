package containers

import "sync"

type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
}

type stack[T any] struct {
	sync.Mutex
	list []T
	head int
}

func NewStack[T any](init ...T) *stack[T] {
	return &stack[T]{
		list: init,
		head: len(init),
	}
}

func (s *stack[T]) Push(value T) {
	s.Lock()
	defer s.Unlock()
	s.list = append(s.list, value)
	s.head++
}

func (s *stack[T]) Pop() (T, bool) {
	s.Lock()
	defer s.Unlock()
	if s.head == 0 {
		var t T
		return t, false
	}
	s.head--
	return s.list[s.head], true
}
