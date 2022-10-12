package analyzer

import "fmt"

// IntLiteral represents a constant signed integer value.
type IntLiteral struct {
	locatable
	value int64
}

// UIntLiteral represents a constant unsigned itneger value.
type UIntLiteral struct {
	locatable
	value uint64
}

// FloatLiteral represents a constant floating point value.
type FloatLiteral struct {
	locatable
	value float64
}

// StringLiteral represents a constant text value.
type StringLiteral struct {
	locatable
	value string
}

// ToString outputs the data in the argument as a string.
func (literal IntLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg int ", literal.value, "\n"))
	return
}

// What returns the type of the argument
func (literal IntLiteral) What () (what Type) {
	what.actual = &PrimitiveI64
	what.length = 1
	return
}

// canBePassedAs returns true if this literal can be implicitly cast to the
// specified type, and false if it can't.
func (literal IntLiteral) canBePassedAs (what Type) (allowed bool) {
	// must be a singlular value
	if !what.isSingular() { return }
	
	// can be passed to types that are signed numbers at a primitive level.
	primitive := what.underlyingPrimitive()
	switch primitive {
	case
		&PrimitiveF64,
		&PrimitiveF32,
		&PrimitiveI64,
		&PrimitiveI32,
		&PrimitiveI16,
		&PrimitiveI8,
		&PrimitiveInt:

		allowed = true
	}
	return
}

// ToString outputs the data in the argument as a string.
func (literal UIntLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg uint ", literal.value, "\n"))
	return
}

// What returns the type of the argument
func (literal UIntLiteral) What () (what Type) {
	what.actual = &PrimitiveU64
	what.length = 1
	return
}

// canBePassedAs returns true if this literal can be implicitly cast to the
// specified type, and false if it can't.
func (literal UIntLiteral) canBePassedAs (what Type) (allowed bool) {
	// must be a singlular value
	if !what.isSingular() { return }
	
	// can be passed to types that are numbers at a primitive level.
	primitive := what.underlyingPrimitive()
	switch primitive {
	case
		&PrimitiveF64,
		&PrimitiveF32,
		&PrimitiveI64,
		&PrimitiveI32,
		&PrimitiveI16,
		&PrimitiveI8,
		&PrimitiveInt,
		&PrimitiveU64,
		&PrimitiveU32,
		&PrimitiveU16,
		&PrimitiveU8,
		&PrimitiveUInt:

		allowed = true
	}
	return
}

// What returns the type of the argument
func (literal FloatLiteral) What () (what Type) {
	what.actual = &PrimitiveF64
	what.length = 1
	return
}

// ToString outputs the data in the argument as a string.
func (literal FloatLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg float ", literal.value, "\n"))
	return
}

// canBePassedAs returns true if this literal can be implicitly cast to the
// specified type, and false if it can't.
func (literal FloatLiteral) canBePassedAs (what Type) (allowed bool) {
	// must be a singlular value
	if !what.isSingular() { return }

	// can be passed to types that are floats at a primitive level.
	primitive := what.underlyingPrimitive()
	switch primitive {
	case
		&PrimitiveF64,
		&PrimitiveF32:

		allowed = true
	}
	return
}

// What returns the type of the argument
func (literal StringLiteral) What () (what Type) {
	what.actual = &BuiltInString
	what.length = 1
	return
}

// ToString outputs the data in the argument as a string.
func (literal StringLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg string '", literal.value, "'\n"))
	return
}

// canBePassedAs returns true if this literal can be implicitly cast to the
// specified type, and false if it can't.
func (literal StringLiteral) canBePassedAs (what Type) (allowed bool) {
	// can be passed to types that are numbers at a primitive level, or
	// types that can be reduced to a variable array pointing to numbers at
	// a primitive level.

	// we don't check the length of what, becasue when setting a static
	// array to a string literal, excess data will be cut off (and if it is
	// shorter, the excess space will be filled with zeros).
	
	reduced, worked := what.reduce()
	if worked {
		if !what.isSingular() { return }
		if reduced.kind != TypeKindVariableArray { return }
		what = reduced
	}

	primitive := what.underlyingPrimitive()
	switch primitive {
	case
		&PrimitiveF64,
		&PrimitiveF32,
		&PrimitiveI64,
		&PrimitiveI32,
		&PrimitiveI16,
		&PrimitiveI8,
		&PrimitiveInt,
		&PrimitiveU64,
		&PrimitiveU32,
		&PrimitiveU16,
		&PrimitiveU8,
		&PrimitiveUInt:

		allowed = true
	}
	return
}
