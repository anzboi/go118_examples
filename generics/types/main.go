package main

func AppendInts(ints Set[int], elems ...int) {
	for i := range elems {
		ints.Add(i)
	}
}

func Append[T comparable](s Set[T], elems ...T) {
	for _, elem := range elems {
		s.Add(elem)
	}
}

func main() {
	var ints Set[int]
	ints = set[int]{}
	ints.Add(1)
	Append(ints, 1, 2, 3)
}
