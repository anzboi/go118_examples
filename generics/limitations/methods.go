package main

import "fmt"

type GenericMethods interface {
	Ident[T any]() T
}

type MyStruct struct {
	field string
}

func (m MyStruct) AddField[T any](val T) string {
	return fmt.Sprintf("%v%s", val, m.field)
}
