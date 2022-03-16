package main

import (
	"fmt"

	"github.com/anzboi/go118-examples/generics/stream"
)

func main() {
	arr := []int{1, 2, 2, 3, 4, 7}
	firstOdd, ok := stream.NewStreamOf(arr...).
		Filter(func(i int) bool { return i%2 == 1 }).
		First()

	if !ok {
		panic("No odd numbers?!?")
	}
	fmt.Println(firstOdd)
}
