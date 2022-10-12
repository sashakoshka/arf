package analyzer

// This is a global, cannonical list of primitive and built-in types.

// PrimitiveF32 is a 32 bit floating point primitive.
var PrimitiveF32  = createPrimitive("F32",  Type {})

// PrimitiveF64 is a 64 bit floating point primitive.
var PrimitiveF64  = createPrimitive("F64",  Type {})

// PrimitiveInt is a signed integer word primitive.
var PrimitiveInt  = createPrimitive("Int",  Type {})

// PrimitiveUInt is an unsigned integer word primitive.
var PrimitiveUInt = createPrimitive("UInt", Type {})

// PrimitiveI8 is a signed 8 bit integer primitive.
var PrimitiveI8   = createPrimitive("I8",   Type {})

// PrimitiveI16 is a signed 16 bit integer primitive.
var PrimitiveI16  = createPrimitive("I16",  Type {})

// PrimitiveI32 is a signed 32 bit integer primitive.
var PrimitiveI32  = createPrimitive("I32",  Type {})

// PrimitiveI64 is a signed 64 bit integer primitive.
var PrimitiveI64  = createPrimitive("I64",  Type {})

// PrimitiveI8 is an unsigned 8 bit integer primitive.
var PrimitiveU8   = createPrimitive("U8",   Type {})

// PrimitiveI16 is an unsigned 16 bit integer primitive.
var PrimitiveU16  = createPrimitive("U16",  Type {})

// PrimitiveI32 is an unsigned 32 bit integer primitive.
var PrimitiveU32  = createPrimitive("U32",  Type {})

// PrimitiveI64 is an unsigned 64 bit integer primitive.
var PrimitiveU64  = createPrimitive("U64",  Type {})

// PrimitiveObj is a blank object primitive.
var PrimitiveObj  = createPrimitive("Obj",  Type {})

// PrimitiveFace is a blank interface primitive. It accepts any value.
var PrimitiveFace = createPrimitive("Face", Type {})

// PrimitiveFunc is a blank function interface primitive. It is useless.
var PrimitiveFunc = createPrimitive("Func", Type {})

// BuiltInString is a built in string type. It is a dynamic array of UTF-32
// codepoints.
var BuiltInString = createPrimitive("String", Type {
	actual: &PrimitiveU32,
	kind:   TypeKindVariableArray,
})

// createPrimitive provides a quick way to construct a primitive for the above
// list.
func createPrimitive (name string, inherits Type) (primitive TypeSection) {
	primitive.where = locator { name: name }
	primitive.what  = inherits
	return
}
