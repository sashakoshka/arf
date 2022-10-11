package analyzer

import "fmt"

type IntLiteral    int64
type UIntLiteral   uint64
type FloatLiteral  float64
type StringLiteral string
type RuneLiteral   rune

// ToString outputs the data in the argument as a string.
func (literal IntLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg int ", literal, "\n"))
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
	output += doIndent(indent, fmt.Sprint("arg uint ", literal, "\n"))
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

// ToString outputs the data in the argument as a string.
func (literal FloatLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg float ", literal, "\n"))
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

// ToString outputs the data in the argument as a string.
func (literal StringLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg string \"", literal, "\"\n"))
	return
}

// canBePassedAs returns true if this literal can be implicitly cast to the
// specified type, and false if it can't.
func (literal StringLiteral) canBePassedAs (what Type) (allowed bool) {
	// can be passed to types that are numbers at a primitive level, or
	// types that can be reduced to a variable array pointing to numbers at
	// a primitive level.

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

// ToString outputs the data in the argument as a string.
func (literal RuneLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg rune '", literal, "'\n"))
	return
}

