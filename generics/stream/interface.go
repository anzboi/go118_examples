package main

type Stream[T any] interface {
	ForEach(func(T))
	Filter(func(T) bool) Stream[T]
	Count() int
	AsList() []T
	First() (T, bool)
}
