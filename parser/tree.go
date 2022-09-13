package parser

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

// SyntaxTree represents an abstract syntax tree. It covers an entire module. It
// can be expected to be syntactically correct, but it might not be semantically
// correct (because it has not been analyzed yet.)
type SyntaxTree struct {
	license string
	author  string

	requires map[string] string
	sections map[string] Section
}

// SectionKind differentiates Section interfaces.
type SectionKind int

const (
	SectionKindType = iota
	SectionKindEnum
	SectionKindFace
	SectionKindData
	SectionKindFunc
)

// Section can be any kind of section. You can find out what type of section it
// is with the Kind method.
type Section interface {
	Location   () (location file.Location)
	Kind       () (kind SectionKind)
	Permission () (permission types.Permission)
	Name       () (name string)
	NewError   (message string, kind infoerr.ErrorKind) (err error)
	ToString   (indent int) (output string)
}

// Identifier represents a chain of names separated by a dot.
type Identifier struct {
	locatable
	trail []string
}

// TypeKind represents what kind of type a type is.
type TypeKind int

const (
	// TypeKindBasic means its a normal type and inherits from something.
	// Basic types can define new members on their parent types.
	TypeKindBasic TypeKind = iota

	// TypeKindPointer means it's a pointer.
	TypeKindPointer

	// TypeKindVariableArray means it's an array of variable length.
	TypeKindVariableArray
)

// TypeMember represents a member variable of a type specifier.
type TypeMember struct {
	locatable
	nameable
	typeable
	permissionable
	
	bitWidth uint64
}

// Type represents a type specifier
type Type struct {
	locatable

	mutable bool
	kind TypeKind
	length uint64

	// only applicable for basic.
        name Identifier

	// not applicable for basic.
	points *Type

	// if non-nil, this type defines new members.
	members      []TypeMember

	// the default value of the type.
	defaultValue Argument
}

// Declaration represents a variable declaration.
type Declaration struct {
	locatable
	nameable
	typeable
}

// ObjectDefaultValues represents a list of object member initialization
// attributes.
type ObjectDefaultValues struct {
	locatable
	attributes map[string] Argument
}

// ArrayDefaultValues represents a list of elements initializing an array.
type ArrayDefaultValues struct {
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

	// (.name <value>)
	// (.name <value> .name (.name <value))
	ArgumentKindObjectDefaultValues

	// <4 32 98 5>
	ArgumentKindArrayDefaultValues

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
	typeable
	permissionable

	external bool
}

// TypeSection represents a type definition.
type TypeSection struct {
	locatable
	nameable
	typeable
	permissionable
}

// EnumMember represents a member of an enum section.
type EnumMember struct {
	locatable
	nameable
	valuable
}

// EnumSection represents an enumerated type section.
type EnumSection struct {
	locatable
	nameable
	typeable
	permissionable

	members []EnumMember
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
	permissionable
	inherits Identifier
	
	behaviors map[string] FaceBehavior
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
	returnees []Argument

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
}

// FuncSection represents a function section.
type FuncSection struct {
	locatable
	nameable
	permissionable
	
	receiver *Declaration
	inputs   []Declaration
	outputs  []FuncOutput
	root     Block

	external bool
}
