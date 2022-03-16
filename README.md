# Go 1.18

Demonstrations of new features in go 1.18.

- [Release notes](https://go.dev/doc/go1.18)
- [Release blog post](https://go.dev/blog/go1.18)
- [Download](https://go.dev/dl/)

## [Generics](./generics)

Golang has generics now!

- The [new generics tutorial]() contains the basics.
- See our examples under the [generics](./generics) directory. We cover our opinions on them, when to use them, and common examples and pitfalls.

## Fuzzer

TODO

## Build info

TODO

# Speach notes

- Generics overview/crash course
- Type parameters/declaration
- Generic Interfaces
- Constraints
- Limitations

## Crash course

- Generics lets you write code that operates on many types.
- The type is fixed at compile time wherever the code is used.

### Generics vs interfaces

- Solves a similar requirement of a language - polymorphism
- Solves it in a very different way.
- Interfaces - same behaviour expected FROM an objects
- Generics - same code acting ON an object.
- Not mutually exclusive (container types, streams/functional programming)

### Declaration/Instantiation

- Functions, structs, interfaces
- Instantiation

### Constraints

- Interfaces
- Unions

### Limitations

The Most severe limitations can be understood by realizing one fact

> It is hard (if not impossible) to tell when an instantiated type implements a generic interface without telling the compiler.

- Don't expect types to automatically implement generic interfaces. Set up a variable with the type parameter set.
- No generic methods.

Some other nuisances.

- No operator overloading makes it harder to write some generic code (eg: how do you check for equality between two objects?).
- Switching on union types is awkward (but its a code smell anyway).
