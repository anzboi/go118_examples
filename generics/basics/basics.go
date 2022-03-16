package main

import "fmt"

// Generic code, type parameters in [] after the function name
func AsList[T any](values ...T) []T {
	// simply returns the variadic as a list
	return values
}

func main() {
	list := AsList(1, 2, 3) // returns a []int
	fmt.Println(list)
}
