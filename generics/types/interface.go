package main

type Set[T comparable] interface {
	Has(T) bool
	Add(T)
}
