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
)

// Type represents a description of a type. It must eventually point to a
// TypeSection.
type Type struct {
	// one of these must be nil.
	actual *TypeSection
	points *Type

	mutable bool
	kind TypeKind

	primitiveCache *TypeSection

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
	}

	if what.points != nil {
		output += " {\n"
		output += what.points.ToString(indent + 1)
		output += doIndent(indent, "}")
	}

	if what.actual != nil {
		output += " " + what.actual.Name()
	}

	output += "\n"
	return
}

// underlyingPrimitive returns the primitive that this type eventually inherits
// from.
func (what Type) underlyingPrimitive () (underlying *TypeSection) {
	// if we have already done this operation, return the cahced result.
	if what.primitiveCache != nil {
		underlying = what.primitiveCache
		return
	}

	if what.kind != TypeKindBasic {
		// if we point to something, return nil because there is no void
		// pointer bullshit in this language
		return
	}

	actual := what.actual
	switch actual {
	case
		&PrimitiveF32,
		&PrimitiveF64,
		&PrimitiveFunc,
		&PrimitiveFace,
		&PrimitiveObj,
		&PrimitiveU64,
		&PrimitiveU32,
		&PrimitiveU16,
		&PrimitiveU8,
		&PrimitiveI64,
		&PrimitiveI32,
		&PrimitiveI16,
		&PrimitiveI8,
		&PrimitiveUInt,
		&PrimitiveInt:

		underlying = actual
		return
	
	case nil:
		panic("invalid state: Type.actual is nil")

	default:
		// if none of the primitives matched, recurse.
		underlying = actual.what.underlyingPrimitive()
		return
	}
}

// reduce ascends up the inheritence chain and gets the first type it finds that
// isn't basic. If the type has a clear path of inheritence to a simple
// primitive, there will be no non-basic types in the chain and this method will
// return false for reducible. If the type this method is called on is not
// basic, it itself is returned.
func (what Type) reduce () (reduced Type, reducible bool) {
	reducible = true

	// returns itself if it is not basic (cannot be reduced further)
	if what.kind != TypeKindBasic {
		reduced = what
		return
	}

	// if we can't recurse, return false for reducible
	if what.actual == nil {
		reducible = false
		return
	}

	// otherwise, recurse
	reduced, reducible = what.actual.what.reduce()
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
		var actual Section
		actual,
		bitten,
		err = analyzer.fetchSectionFromIdentifier(inputType.Name())

		outputType.actual = actual.(*TypeSection)
		// TODO: produce an error if this doesnt work
		
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
