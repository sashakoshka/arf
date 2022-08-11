# ARF

The ARF programming language.

This is still under development and does not compile things yet. Once complete,
it will serve as a temporary compiler that will be used to write a new one using
the language itself.

The old repository can be found [here](https://github.com/sashakoshka/arf-old).

ARF is a low level language with a focus on organization, modularization, and
code clarity. Behind it's avant-garde syntax, its basically just a more refined
version of C.

A directory of ARF files is called a module, and modules will compile to object
files (one per module) using C as an intermediate language (maybe LLVM IR in the
future).

## Design aspects

These are some design goals that I have followed/am following:

- The standard library will be fully optional, and decoupled from the language
- The language itself must be extremely simple
- Language features must be immutable (no reflection or operator overloading)
- Data must be immutable by default
- Memory not on the stack must be allocated and freed manually
- Language syntax must have zero ambiguity

## Planned features

- Type definition through inheritence
- Struct member functions
- Go-style interfaces
- Generics
- A standard library (that can be dynamically linked)

## Checklist

- [X] File reader
- [X] File -> tokens
- [ ] Tokens -> syntax tree
- [ ] Syntax tree -> semantic tree
- [ ] Semantic tree -> C -> object file
- [ ] Figure out HOW to implement generics
- [ ] Create a standard library
