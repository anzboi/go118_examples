package main

import (
	"constraints"
	"fmt"
)

type Vec2D[F constraints.Float] struct {
	x, y F
}

func Foo[T, S any](arg1 T) S {
	var s S
	return s
}

func main() {
	s := Foo[int, string](1)
	fmt.Println(s)
}
