package analyzer

// This is a global, cannonical list of primitive and built-in types.

var PrimitiveInt  = TypeSection { name: "Int" }
var PrimitiveUInt = TypeSection { name: "UInt" }
var PrimitiveI8   = TypeSection { name: "I8" }
var PrimitiveI16  = TypeSection { name: "I16" }
var PrimitiveI32  = TypeSection { name: "I32" }
var PrimitiveI64  = TypeSection { name: "I64" }
var PrimitiveU8   = TypeSection { name: "U8" }
var PrimitiveU16  = TypeSection { name: "U16" }
var PrimitiveU32  = TypeSection { name: "U32" }
var PrimitiveU64  = TypeSection { name: "U64" }

var PrimitiveObjt = TypeSection { name: "Objt" }
var PrimitiveFace = TypeSection { name: "Face" }

var BuiltInString = TypeSection {
	inherits: Type {
		actual: PrimitiveU8,
		kind:   TypeKindVariableArray,
	},
}
