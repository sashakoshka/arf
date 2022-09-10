package analyzer

// This is a global, cannonical list of primitive and built-in types.

var PrimitiveInt  = createPrimitive("Int",  Type {})
var PrimitiveUInt = createPrimitive("UInt", Type {})
var PrimitiveI8   = createPrimitive("I8",   Type {})
var PrimitiveI16  = createPrimitive("I16",  Type {})
var PrimitiveI32  = createPrimitive("I32",  Type {})
var PrimitiveI64  = createPrimitive("I64",  Type {})
var PrimitiveU8   = createPrimitive("U8",   Type {})
var PrimitiveU16  = createPrimitive("U16",  Type {})
var PrimitiveU32  = createPrimitive("U32",  Type {})
var PrimitiveU64  = createPrimitive("U64",  Type {})
var PrimitiveObjt = createPrimitive("Objt", Type {})
var PrimitiveFace = createPrimitive("Face", Type {})

var BuiltInString = createPrimitive("String", Type {
	actual: PrimitiveU8,
	kind:   TypeKindVariableArray,
})

// createPrimitive provides a quick way to construct a primitive for the above
// list.
func createPrimitive (name string, inherits Type) (primitive TypeSection) {
	primitive.where = locator { name: name }
	primitive.inherits = inherits
	return
}
