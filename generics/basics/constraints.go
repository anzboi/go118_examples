package main

import "fmt"

// Type parameter must be a string or int
type MyConstraint interface {
	int | string
}

func DoSomething[T MyConstraint](input T) {
	fmt.Println(input)
}

func example() {
	DoSomething("123")
}

// Type parameter must implement fmt.Stringer and alias string or int.
type MyAdvancedConstraint interface {
	fmt.Stringer
	~int | ~string
}

type IP string

func (ip IP) String() string {
	return string(ip)
}

func DoAnother[T MyAdvancedConstraint](input T) {
	fmt.Println(input.String())
}
