package analyzer

// import "git.tebibyte.media/arf/arf/types"
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
	outputType.length  = inputType.Length()
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
	} else {
		var bitten parser.Identifier
		outputType.actual,
		bitten,
		err = analyzer.fetchSectionFromIdentifier(inputType.Name())
		
		if bitten.Length() > 0 {
			err = bitten.NewError(
				"cannot use member selection in this context",
				infoerr.ErrorKindError)
			return
		}
	}
	
	// TODO
	
	return
}
