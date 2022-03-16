# Generics

Golang has generics now!

- [basic tutorial](https://go.dev/doc/tutorial/generics)
- [generics language proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
- [Why generics (go blog)](https://go.dev/blog/why-generics)

For the rest of this part of the repo, we discuss where generics fits in the go language, and give examples of some common use cases and pitfalls.

## Contents

1. [Where do generics fit](#where-do-generics-fit)
2. [Key Concepts](#key-concepts)
3. [Generic Code](#generic-code)
2. [Limitations](#limitations)
2. [Common pitfalls](#common-pitfalls)

## Where do generics fit

Generics are somewhat related to interfaces. They both implement a form of polymorphism to encourage code reuse, reduce code duplication. However the way generics and interfaces provide polymorphism is very different.

### Your Program vs Your Code

To understand the theory about generics we first must talk about the difference between a program at an abstract level, and the code you write to implement it.

The program consists of ideas such as objects, functions, function calls, pointers, control flows, slices, goroutines, channels and more. These are all abstractions on what your program is doing and what information its storing as it runs on your machine.

Code is distinctly different. Your code contains variables, object definitions, function definitions, operators, comma-separated argument lists. Obviously some ideas map between code and program pretty cleanly, but always keep in mind that your code is a description of your program. This brings us to one key insight to understanding generics.

> There is no such thing as a generic program, only generic code.

### Generic code

Generic code a typed programming language arises almost immediately when you need to write two, almost identicle functions that operate on different types.

```go
func SortInts(arr []int) {
    ...
    if arr[j] < arr[i] {
        arr[i], arr[j] = arr[j], arr[i]
    }
    ...
}

func SortStrings(arr []string) {
    ...
    if arr[j] < arr[i] {
        arr[i], arr[j] = arr[j], arr[i]
    }
    ...
}
```

You will notice the code is (almost) exactly the same, but the program is not. The code has a new function name, and input type, thats it. The program however, can be quite different. String comparison is much more involved than int comparison, and assigning the variables involves copying an array of bytes, instead of a single int.

Generic code allows us to reuse the same code, on a generic type.

```go
func Sort[T constraint](arr []T) {
    // generic code
}
```

It is now the compilers job to fill in the details when generic code is used. If part of the program calls `Sort([]string)`, then the compiler _"creates"_ a function using the generic code and setting `T=string`. Your program does not contain a generic `Sort` function, only a `Sort[string]` function, and `Sort[int]` if that is needed as well.

We can do the same thing with other pieces of code, namely structs and interfaces. The compiler will create real structs and interfaces out of the generic ones at compile time as they are needed.

```go
type Vec2D[T numeric] struct {
    x, y T
}

type MyInterface[T any] interface {
    Method(T) string
}
```

DISCLAIMER: What we have described here is a useful framework to think about generics. The actual implementation may differ. The [design proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#efficiency) does not specify how to implement generics, but provides a few options.

### So where does generic code fit?

Generics fit wherever you need code to operate on multiple types. Maybe your function has a couple of inputs where the types need to match up in some sense, or your output needs to match the input. Traditionally you would use an interface, and expect the caller to cast things as necessary, or simply defined as many functions as there are combinations of inputs.

An example adapted from one of our repos. We have a need to search a list of objects for an object that matches a key, and return the match. Previously we would write a function per object array type that we need to search.

```go
// all our DB objects implement this.
type Keyer interface {
    GetName() string
}

func FindObject1(key string, objects []*Object1) *Object1 {
    for _, obj := range objects {
        if obj.GetName() == key {
            return obj
        }
    }
    return nil
}
func FindObject2(key string, objects []*Object2) *Object2 {
    for _, obj := range objects {
        if obj.GetName() == key {
            return obj
        }
    }
    return nil
}
```

There are numerous tools out there that try to auto generate code like this. Another approach involves using interfaces that need to be implemented by slice type aliases (eg: the std `sort` package).

Instead, we can now define a generic function. We add a bool arg since we cannot assume T will always be a pointer type.

```go
func FindObject[T Keyer](key string, objects []T) (T, bool) {
    for _, obj := range objects {
        if obj.GetName() == key {
            return obj, true
        }
    }
    var zero T // easiest way to instantiate the zero value for the return
    return zero, false
}
```

This is a huge improvement. This will now operate on any object array we may need in the future, and compile to code that is (almost) equivalent to each individual function. We simply need to make sure each new type added implements `Keyer`.

You will notice we did need to make one change though, the line `var zero T`. This is because many different types have wildly different ways to instantiate them. You assign a number to numerics. Slices look like `[]type{}`. The code `var {variable} Type` is a universal way to instantiate an empty object.

> In some sense, `var {variable} Type` is generic enough to be used generically.

## Key concepts

Here we go over the key concepts for how generics are implemented and used in Go.

### Type parameters

Generic code uses type parameters as a placeholder to actual types when the generic code is used. This is 100% analogous to how function arguments are placeholder for actual values in a function call.

```go
func MyGenericFunc[T any](arg1 T, arg2 []T) T {
    // generic code
}

// generics can have multiple type parameters
func MyGenericFunc[T any, S any](arg1 T, arg2 []T) S {
    // generic code
}

// you can group type parameters with the same constraints
func MyGenericFunc[T, S any](arg1 T, arg2 []T) S {
    // generic code
}
```

### Constraints

Constraints limit what types your generic code can operate on. Constraints are interfaces, and interfaces have recieved some upgrades to support some richer constraint types.

Go also has a new keyword, `any` that is simply an alias of `interface{}`. It can and should be used in its place. As a constraint, `any` says that any a type parameter could be any type.

NOTE: The main change to interfaces is that they can now be declared with a types list as well as a method list. Interfaces with a non-empty type list cannot be used as interfaces in the usual sense and can only be used as generic type constraints.

```go
// This interface is effectively a union constraint of string|int|MyStruct
type MyConstraint interface {
    string | int | MyStruct
}

Func DoSomething[T MyConstraint](item T) {} // item must be either a string, int or MyStruct.

Func DoSomething[T any](item T) {} // item can have any type

Func DoSomething(item any) {} // any can be used in place of interface{}
```

### Alias constraints

Golang supports constraints where the type aliases an underlying type

```go
// PrintInt accepts types where the underlying type is int
func PrintInt[T ~int](val T) {
    fmt.Printf("%d", val)
}

type MyInt int
one := MyInt(1)
PrintInt(one)
```

### Mixed

A more complicated example involving methods and types in a single constraint.

```go
// This constraint requires that types are aliases of int or string, AND implement fmt.Stringer
type MyConstraint interface {
    fmt.Stringer
    ~int | ~string
}
```

### Instantiation

Recall that there is no generic programs, only generic code. Conversion from generic code to real code happens at compile time, and this means the compiler needs to be able to resolve generic types to real types before it can convert it into the final program. In other words...

> The compiler needs to be able to instantiate real variables, functions, etc... Not generic ones.

```go
// generic struct
type Foo[T any] {
    value T
}

// This is Invalid, type parameter is not present when constructing Foo.
var f Foo
f = Foo[string]{}

// Better, we are telling the compiler which type to substitute
var g Foo[string]

// We can declare and assign in one go
h := Foo[string]{}

// Go can also infer the type parameter from input.
j := Foo{value: "Hello"}
```

This includes function calls and interface variables.

```go
type Bar[T any] interface {}
func Baz[T any](value T) string {}

var b Bar[int]
str := Baz[string]("hello")
```

### Inference

The go compiler makes efforts to infer type parameters where possible. Most of the time, you will not have to specify type parameters when using generic code.

```go
type Foo[T any] struct {
    value T
}
func NewFoo[T any](value T) *Foo[T] {
    return &Foo{value: value}
}

f1 := &Foo{value: "hello world"} // Infers T = string from input
f2 := NewFoo("hello world") // Same
```


## Limitations

There are a number of limitations with the current generic implementation in go 1.18. For a more complete list, see the [omissions section](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#omissions) of the design proposal.

### Type Inference and Methods

A class of limitations can be traced back to one single fact.

> It is hard to know when a function signature can match the type signature of a generic function given a suitable combination of type parameters.

Consider the following example.

```go
type Cache[T any] interface {
    Put(T)
    Value() T
}

type IntCache struct {
    sync.RWMutex
    value int
}

func (i *IntCache) Value() int {
    i.RLock()
    defer i.RUnlock()
    return i.value
}

func (i *IntCache) Put(num int) {
    i.Lock()
    defer i.Unlock()
    i.value = num
}
```

Here we have a generic cache interface that stores and fetches a value, and a concrete implementation of the cache where `T = int`. The trouble is that its not easy in general for the compiler to detect the fact that `*IntCache` implements `Cache[int]` (how does it know in advance that `T = int` works).

This limitation has a few effects. Firstly, no generic methods.

```go
// No generic methods
type Foo interface {
    Method[T any]() // Invalid, methods cannot have type parameters
}

// You cannot put generic methods on structs either
type Bar struct {}
func (b *Bar) Baz[T any]() {}
```

Fun fact: its possible to create a scenario where the go program needs _Just In Time_ compilation (Or the compiler needs to walk the entire call graph) to resolve type parameters at compile time.

You also sometimes have to give the compiler the type parameters if it cannot resolve them.

```go
type Foo[T any] interface {
    Get() T
}
type MyString string // MyString implements Foo[string]
func (m MyString) Get() string { return string(m) }

func Bar[T any](input Foo[T]) {}
input := MyString("Hello World")

// invalid, cannot infer type parameter T
// Even though MyString implements Foo[string], the compiler cannot figure out that T=string is a suitable substitution.
Bar(input)

// This works
Bar[string](input)
```

### Lack of type switching

Golang generics do not support type switching for union types. If you think back to what generics are, this makes sense. Generic code is code that operates on multiple possible types. If you fill in your type parameter, it **_usually_** does not make sense to do a type switch.

> Attempting to type switch on a generic type is a code smell.

This does not mean attempting to type switch on a generic is wrong, it should just be considered a warning that maybe your design can be improved.

However, sometimes, you really do want to do a type switch. Consider the following general `Join` function.

```go
func Join[T any](elems []T, joiner string) string {
    b := strings.Builder{}
    for n, elem := range elem {
        if stringer, ok := elem.(fmt.Stringer); ok {
            b.WriteString(stringer.String())
        } else {
            fmt.Fprintf(&b, "%v", elem)
        }
        if n < len(elems) {
            b.WriteString(joiner)
        }
    }
}
```

We could improve the efficiency of this code for simple types by delegating the join function to type-specific functions IFF we could switch on the type parameter.

```go
func Join[T any](elems []T, joiner string) string {
    switch T {
    case ~int:
        // int specific function (eg: use %d in the format string)
    case ~rune:
    ...
    default:
        // as above
    }
}
```

For another example, see the [order repo](https://github.com/jamesrom/order). This repo exists because you cannot type switch on a type parameter to change ordering functionality.

### No operator overloading

Go doesn't provide operator overloading. This has a negative effect on generic code, since not all code is applicable on various types. For example, there is no way to make a struct implement a `<` operator.

```go
// ordered are things that support '<'
func Sort[T ordered](arr []T) []T {
    ...
    if arr[j] < arr[i] {
        // swap
    }
}

// There is no way to make Foo support '<', so you cannot Sort an array of Foo using the above sort function.
type Foo struct {
    ...
}
```

You can define and implement a `Less` function, but now your generic code needs to call that function, meaning you can't use `<` on types that DO support it.

```go
type Ordered[T any] interface { Less(T) bool }
type Foo struct { val string }
func (f *Foo) Less(g *Foo) bool { return f.val < g.val }

// ordered are things that support '<'
func Sort[T ordered](arr []T) []T {
    ...
    if arr[j].Less(arr[i]) {
        // swap
    }
}

Sort([]*Foo{{"b"}, {"a"}}) // this works now
Sort([]int{3, 2, 1}) // invalid: you've lost the ability to sort an array of int.
```

## Common Pitfalls

TODO