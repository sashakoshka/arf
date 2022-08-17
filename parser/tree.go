package parser

import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/types"

// SyntaxTree represents an abstract syntax tree. It covers an entire module. It
// can be expected to be syntactically correct, but it might not be semantically
// correct (because it has not been analyzed yet.)
type SyntaxTree struct {
	license string
	author  string

	requires     []string
	dataSections map[string] *DataSection
}

// Identifier represents a chain of arguments separated by a dot.
type Identifier struct {
	location file.Location
	trail    []string
}

// TypeKind represents what kind of type a type is
type TypeKind int

const (
	// TypeKindBasic either means it's a primitive, or it inherits from
	// something.
	TypeKindBasic TypeKind = iota

	// TypeKindPointer means it's a pointer
	TypeKindPointer

	// TypeKindArray means it's an array.
	TypeKindArray
)

// Type represents a type specifier
type Type struct {
	location file.Location

	mutable bool
	kind TypeKind

	// only applicable for arrays. a value of zero means it has an
	// undefined/dynamic length.
	length uint64

	// only applicable for basic.
        name Identifier

	// not applicable for basic.
	points *Type
}

// Declaration represents a variable declaration.
type Declaration struct {
	location file.Location
	name     string
	what     Type
}

// ObjectInitializationValues represents a list of object member initialization
// attributes.
type ObjectInitializationValues struct {
	location   file.Location
	attributes map[string] Argument
}

// ArrayInitializationValues represents a list of attributes initializing an
// array.
type ArrayInitializationValues struct {
	location file.Location
	values   []Argument
}

// Phrase represents a function call or operator. In ARF they are the same
// syntactical concept.
type Phrase struct {
	location  file.Location
	command   Argument
	arguments []Argument
	returnsTo []Argument
}

// ArgumentKind specifies the type of thing the value of an argument should be
// cast to.
type ArgumentKind int

const (
	// [name argument]
	// [name argument argument]
	// etc...
	ArgumentKindPhrase ArgumentKind = iota

	// {name}
	ArgumentKindDereference
	
	// {name 23}
	ArgumentKindSubscript

	// .name value
	// but like, a lot of them
	ArgumentKindObjectInitializationValues

	// value value...
	ArgumentKindArrayInitializationValues

	// name.name
	// name.name.name
	// etc...
	ArgumentKindIdentifier

	// name:Type
	// name:{Type}
	// name:{Type ..}
	// name:{Type 23}
	// etc...
	ArgumentKindDeclaration

	// -1337
	ArgumentKindInt

	// 1337
	ArgumentKindUInt

	// 0.44
	ArgumentKindFloat

	// "hello world"
	ArgumentKindString

	// 'S'
	ArgumentKindRune

	// + - * / etc...
	// this is only used as a phrase command
	ArgumentKindOperator
)

// Argument represents a value that can be placed anywhere a value goes. This
// allows things like phrases being arguments to other phrases.
type Argument struct {
	location file.Location
	kind     ArgumentKind
	value    any
	// TODO: if there is an argument expansion operator its existence should
	// be stored here in a boolean.
}

// DataSection represents a global variable.
type DataSection struct {
	location file.Location
	name     string
	
	what       Type
	value      Argument
	permission types.Permission
}
