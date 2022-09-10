package analyzer

// This is a global, cannonical list of primitive and built-in types.

var PrimitiveInt  = TypeSection { sectionBase: sectionBase { name: "Int" }  }
var PrimitiveUInt = TypeSection { sectionBase: sectionBase { name: "UInt" } }
var PrimitiveI8   = TypeSection { sectionBase: sectionBase { name: "I8 " }  }
var PrimitiveI16  = TypeSection { sectionBase: sectionBase { name: "I16 " } }
var PrimitiveI32  = TypeSection { sectionBase: sectionBase { name: "I32 " } }
var PrimitiveI64  = TypeSection { sectionBase: sectionBase { name: "I64 " } }
var PrimitiveU8   = TypeSection { sectionBase: sectionBase { name: "U8 " }  }
var PrimitiveU16  = TypeSection { sectionBase: sectionBase { name: "U16 " } }
var PrimitiveU32  = TypeSection { sectionBase: sectionBase { name: "U32 " } }
var PrimitiveU64  = TypeSection { sectionBase: sectionBase { name: "U64 " } }

var PrimitiveObjt = TypeSection { sectionBase: sectionBase { name: "Objt" } }
var PrimitiveFace = TypeSection { sectionBase: sectionBase { name: "Face" } }

var BuiltInString = TypeSection {
	inherits: Type {
		actual: PrimitiveU8,
		kind:   TypeKindVariableArray,
	},
}
