package iterator

type Iterator[T any] interface {
	Next() bool
	Value() T
}

type iter[T any] struct {
	elems []T
	ptr   int
}

func NewIterator[T any](elems []T) *iter[T] {
	return &iter[T]{
		elems: elems,
		ptr:   -1,
	}
}

func (i *iter[T]) Next() bool {
	if i.ptr == len(i.elems)-1 {
		return false
	}
	i.ptr++
	return true
}

func (i *iter[T]) Value() T {
	return i.elems[i.ptr]
}
