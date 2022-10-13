package analyzer

import "fmt"
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
)

// Type represents a description of a type. It must eventually point to a
// TypeSection.
type Type struct {
	locatable

	// one of these must be nil.
	actual Section
	points *Type

	mutable bool
	kind TypeKind

	primitiveCache *TypeSection
	singularCache  *bool

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
// from. If the type ends up pointing to something, this returns nil.
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

		underlying = actual.(*TypeSection)
		return
	
	case nil:
		panic (
			"invalid state: Type.actual is nil for " +
			what.Describe() + " " +
			what.locatable.location.Describe())

	default:
		// if none of the primitives matched, recurse.
		switch actual.(type) {
		case *TypeSection:
			underlying =
				actual.(*TypeSection).
				what.underlyingPrimitive()
		// TODO
		// case *FaceSection:
			// TODO: depending on if this is an object interface or
			// a function interface, return either Face or Func.
			// we can assume this because of inheritence rules.
			
		// // case *EnumSection:
			// underlying =
				// actual.(*EnumSection).
				// what.underlyingPrimitive()
		
		default:
			panic (
				"invalid state: type " + what.Describe() +
				"has illegal actual " +
				what.locatable.location.Describe())
		}
		return
	}
}

// isSingular returns whether or not the type is a singular value. this goes
// all the way up the inheritence chain, only stopping when it hits a non-basic
// type because this is about data storage of a value.
func (what Type) isSingular () (singular bool) {
	// if we have already done this operation, return the cahced result.
	if what.singularCache != nil {
		singular = *what.singularCache
		return
	}

	if what.length > 0 {
		singular = true
		return
	}

	// decide whether or not to recurse
	if what.kind != TypeKindBasic { return }
	actual := what.actual
	if actual == nil {
		return
	} else {
		switch actual.(type) {
		case *TypeSection:
			singular = actual.(*TypeSection).what.isSingular()
		// TODO: uncomment this when these sections have been
		// implemented
		// case *FaceSection:
			// singular = true
			
		// case *EnumSection:
			// singular = actual.(*EnumSection).what.isSingular()
		
		default:
			panic (
				"invalid state: type " + what.Describe() +
				"has illegal actual " +
				what.locatable.location.Describe())
		}
	}
	return
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
	switch what.actual.(type) {
	case *TypeSection:
		reduced, reducible = what.actual.(*TypeSection).what.reduce()
		
	// TODO: uncomment this when these sections have been
	// implemented
	// case *FaceSection:
		// singular = true
		
	// case *EnumSection:
		// reduced, reducible = what.actual.(*EnumSection).what.reduce()
	
	default:
		panic (
			"invalid state: type " + what.Describe() +
			"has illegal actual " +
			what.locatable.location.Describe())
	}

	return
}

// analyzeType analyzes a type specifier.
func (analyzer analysisOperation) analyzeType (
	inputType parser.Type,
) (
	outputType Type,
	err error,
) {
	outputType.mutable  = inputType.Mutable()
	outputType.length   = inputType.Length()
	outputType.location = inputType.Location()
	if outputType.length < 1 {
		err = inputType.NewError (
			"cannot specify a length of zero",
			infoerr.ErrorKindError)
		return
	}

	if inputType.Kind() != parser.TypeKindBasic {
		// analyze type this type points to, if it exists
		var points Type
		points, err = analyzer.analyzeType(inputType.Points())
		outputType.points = &points
	} else {
		// analyze the type section this type uses
		var bitten parser.Identifier
		var external bool
		var actual Section
		
		actual,
		external,
		bitten,
		err = analyzer.fetchSectionFromIdentifier(inputType.Name())
		if err != nil { return }

		if actual == nil {
			err = inputType.NewError (
				"this type does not exist",
				infoerr.ErrorKindError)
			return
		}

		actualIsValidSectionKind := false
		switch actual.(type) {
		// TODO: uncomment once these sections are implemented
		case *TypeSection /* , *EnumSection, *FaceSection */:
			actualIsValidSectionKind = true
		}

		if !actualIsValidSectionKind {
			err = inputType.NewError (
				"this must refer to a type, interface, or enum",
				infoerr.ErrorKindError)
			return
		}

		if external && actual.Permission() == types.PermissionPrivate {
			err = bitten.NewError(
				"this type is private, and cannot be used " +
				"outside of its module",
				infoerr.ErrorKindError)
			return
		}
		
		if bitten.Length() > 0 {
			err = bitten.NewError(
				"cannot use member selection in this context",
				infoerr.ErrorKindError)
			return
		}
	}
	
	return
}

// Describe provides a human readable description of the type. The value of this
// should not be computationally analyzed.
func (what Type) Describe () (description string) {
	if what.kind == TypeKindBasic {
		actual := what.actual
		switch actual {
		case &PrimitiveF32:
			description += "F32"
		case &PrimitiveF64:
			description += "F64"
		case &PrimitiveObj:
			description += "Obj"
		case &PrimitiveU64:
			description += "U64"
		case &PrimitiveU32:
			description += "U32"
		case &PrimitiveU16:
			description += "U16"
		case &PrimitiveU8:
			description += "U8"
		case &PrimitiveI64:
			description += "I64"
		case &PrimitiveI32:
			description += "I32"
		case &PrimitiveI16:
			description += "I16"
		case &PrimitiveI8:
			description += "I8"
		case &PrimitiveUInt:
			description += "UInt"
		case &PrimitiveInt:
			description += "Int"
		// case &PrimitiveFunc:
			// description += "Func"
		// case &PrimitiveFace:
			// description += "Face"
		case &BuiltInString:
			description += "String"
		
		case nil:
			description += "NIL-TYPE-ACTUAL"

		default:
			description += actual.ModuleName() + "." + actual.Name()
			return
		}
	} else {
		description += "{"
		description += what.points.Describe()
		description += "}"
	}

	if what.length > 0 {
		description += fmt.Sprint(":", what.length)
	}

	return
}
