package analyzer

import "fmt"
import "git.tebibyte.media/arf/arf/types"

// TypeKind represents what kind of type a type is.
type TypeKind int

const (
	// TypeKindBasic means it's a single value.
	TypeKindBasic TypeKind = iota

	// TypeKindPointer means it's a pointer
	TypeKindPointer

	// TypeKindVariableArray means it's an array of variable length.
	TypeKindVariableArray

	// TypeKindObject means it's a structured type with members.
	TypeKindObject
)

// ObjectMember is a member of an object type. 
type ObjectMember struct {
	name string
	
	// even if there is a private permission in another module, we still
	// need to include it in the semantic analysis because we need to know
	// how many members objects have.
	permission types.Permission
	
	// TODO: create argument type similar to what's in the parser and have
	// a defaultValue member here.
}

// Type represents a description of a type. It must eventually point to a
// TypeSection.
type Type struct {
	actual Section
	points *Type

	mutable bool
	kind TypeKind

	// if this is greater than 1, it means that this is a fixed-length array
	// of whatever the type is. even if the type is a variable length array.
	// because literally why not.
	length uint64

	// this is only applicable for a TypeKindObject where new members are
	// defined.
	// TODO: do not add members from parent type. instead have a member
	// function to discern whether this type contains a particular member,
	// and have it recurse all the way up the family tree. it will be the
	// translator's job to worry about what members are placed where.
	members []ObjectMember
}

// ToString returns all data stored within the type, in string form.
func (what Type) ToString (indent int) (output string) {
	output += fmt.Sprint("") 
	return
}

// TypeSection represents a type definition section.
type TypeSection struct {
	sectionBase
	inherits Type
}

// Kind returns SectionKindType.
func (section TypeSection) Kind () (kind SectionKind) {
	kind = SectionKindType
	return
}

// ToString returns all data stored within the type section, in string form.
func (section TypeSection) ToString (indent int) (output string) {
	return
}
