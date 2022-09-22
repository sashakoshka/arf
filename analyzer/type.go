package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

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

	what Type
}

func (member ObjectMember) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		member.name, " ",
		member.permission.ToString(),
		"\n")
	output += member.what.ToString(indent + 1)
	return
}

// Type represents a description of a type. It must eventually point to a
// TypeSection.
type Type struct {
	// one of these must be nil.
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
	output += doIndent(indent, "type ", what.length)

	if what.mutable {
		output += " mutable"
	}
	
	switch what.kind {
	case TypeKindBasic:
		output += " basic"
	case TypeKindPointer:
		output += " pointer"
	case TypeKindVariableArray:
		output += " variableArray"
	case TypeKindObject:
		output += " object"
	}

	if what.points != nil {
		output += " {\n"
		output += what.points.ToString(indent + 1)
		output += doIndent(indent, "}")
	}

	if what.actual != nil {
		output += what.actual.Name()
	}
	return
}

// analyzeType analyzes a type specifier.
func (analyzer AnalysisOperation) analyzeType (
	inputType parser.Type,
) (
	outputType Type,
	err error,
) {
	outputType.mutable = inputType.Mutable()
	if outputType.length < 1 {
		err = inputType.NewError (
			"cannot specify a length of zero",
			infoerr.ErrorKindError)
		return
	}

	// analyze type this type points to, if it exists
	if inputType.Kind() != parser.TypeKindBasic {
		var points Type
		points, err = analyzer.analyzeType(inputType.Points())
		outputType.points = &points
	}
	
	// TODO
	
	return
}
