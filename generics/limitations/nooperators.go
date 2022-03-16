package limitations

type MyStruct2 struct{}

func Less[T ~int | ~float32 | MyStruct](a, b T) bool {
	return a < b
}
