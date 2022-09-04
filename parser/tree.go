package parser

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"

// SyntaxTree represents an abstract syntax tree. It covers an entire module. It
// can be expected to be syntactically correct, but it might not be semantically
// correct (because it has not been analyzed yet.)
type SyntaxTree struct {
	license string
	author  string

	requires     []string
	typeSections map[string] *TypeSection
	objtSections map[string] *ObjtSection
	enumSections map[string] *EnumSection
	faceSections map[string] *FaceSection
	dataSections map[string] *DataSection
	funcSections map[string] *FuncSection
}

// Identifier represents a chain of arguments separated by a dot.
type Identifier struct {
	locatable
	trail []string
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
	locatable

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
	locatable
	nameable
	what Type
}

// ObjectInitializationValues represents a list of object member initialization
// attributes.
type ObjectInitializationValues struct {
	locatable
	attributes map[string] Argument
}

// ArrayInitializationValues represents a list of attributes initializing an
// array.
type ArrayInitializationValues struct {
	locatable
	values []Argument
}

// ArgumentKind specifies the type of thing the value of an argument should be
// cast to.
type ArgumentKind int

const (
	ArgumentKindNil ArgumentKind = iota
	
	// [name argument]
	// [name argument argument]
	// etc...
	ArgumentKindPhrase

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
	locatable
	kind  ArgumentKind
	value any
	// TODO: if there is an argument expansion operator its existence should
	// be stored here in a boolean.
}

// DataSection represents a global variable.
type DataSection struct {
	locatable
	nameable
	
	what       Type
	permission types.Permission
	value      Argument
}

// TypeSection represents a blind type definition.
type TypeSection struct {
	locatable
	nameable
	
	inherits     Type
	permission   types.Permission
	defaultValue Argument
}

// ObjtMember represents a part of an object type definition.
type ObjtMember struct {
	locatable
	nameable

	what         Type
	bitWidth     uint64
	permission   types.Permission
	defaultValue Argument
}

// ObjtSection represents an object type definition.
type ObjtSection struct {
	locatable
	nameable

	inherits     Identifier
	permission   types.Permission
	members      []ObjtMember
}

// EnumMember represents a member of an enum section.
type EnumMember struct {
	locatable
	nameable
	value Argument
}

// EnumSection represents an enumerated type section.
type EnumSection struct {
	locatable
	nameable
	
	what       Type
	permission types.Permission
	members    []EnumMember
}

// FaceBehavior represents a behavior of an interface section.
type FaceBehavior struct {
	locatable
	nameable

	inputs  []Declaration
	outputs []Declaration
}

// FaceSection represents an interface type section.
type FaceSection struct {
	locatable
	nameable
	inherits Identifier
	
	permission types.Permission
	behaviors  map[string] FaceBehavior
}

// PhraseKind determines what semantic role a phrase plays.
type PhraseKind int

const (
	PhraseKindCall = iota
	PhraseKindCallExternal
	PhraseKindOperator
	PhraseKindAssign
	PhraseKindReference
	PhraseKindDefer
	PhraseKindIf
	PhraseKindElseIf
	PhraseKindElse
	PhraseKindSwitch
	PhraseKindCase
	PhraseKindWhile
	PhraseKindFor
)

// Phrase represents a function call or operator. In ARF they are the same
// syntactical concept.
type Phrase struct {
	location  file.Location
	command   Argument
	arguments []Argument
	returnsTo []Argument

	kind PhraseKind

	// only applicable for control flow phrases
	block Block
}

// Block represents a scoped/indented block of code.
type Block []Phrase

// FuncOutput represents an input a function section. It is unlike an input in
// that it can have a default value.
type FuncOutput struct {
	Declaration
	defaultValue Argument
}

// FuncSection represents a function section.
type FuncSection struct {
	locatable
	nameable
	permission types.Permission
	
	receiver *Declaration
	inputs   []Declaration
	outputs  []FuncOutput
	root     Block

	external bool
}
