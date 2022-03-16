package main

type MyInterface[T any] interface {
	Value() T
}

type MyString string

func (m MyString) Value() string {
	return string(m)
}

func GetValue[T any](impl MyInterface[T]) T {
	return impl.Value()
}

func tryImplement() {
	myValue := MyString("Hello World")

	// compiler cannot infer the type parameter T.
	// MyString implements MyInterface[string], but its not easy to see that 'string'
	// is the correct substitute for T.
	GetValue(myValue)

	// We must instead tell the compiler that T is string
	GetValue[string](myValue)

	// This also works, we are creating a new variable with type MyInterface[string],
	// the compiler is able to infer T is 'string' from the input.
	var impl MyInterface[string] = myValue
	GetValue(impl)
}
