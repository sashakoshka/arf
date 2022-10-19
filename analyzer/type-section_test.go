package analyzer

import "testing"

func TestTypeSection (test *testing.T) {
	checkTree ("../tests/analyzer/typeSection", false,
`typeSection ro ../tests/analyzer/typeSection/required.aBasic
	type 1 basic Int
typeSection ro ../tests/analyzer/typeSection/required.bBird
	type 1 basic Obj
	member rw wing
		type 1 basic Int
		uintLiteral 2
typeSection ro ../tests/analyzer/typeSection.aBasicInt
	type 1 basic Int
	uintLiteral 5
typeSection ro ../tests/analyzer/typeSection.bOnBasicInt
	type 1 basic aBasicInt
typeSection ro ../tests/analyzer/typeSection.cBasicObject
	type 1 basic Obj
	member ro that
		type 1 basic UInt
	member ro this
		type 1 basic Int
typeSection ro ../tests/analyzer/typeSection.dInheritFromOther
	type 1 basic aBasic
typeSection ro ../tests/analyzer/typeSection.eInheritObject
	type 1 basic cBasicObject
	member ro that
		type 1 basic UInt
		uintLiteral 5
typeSection ro ../tests/analyzer/typeSection.fInheritObjectFromOther
	type 1 basic bBird
	member ro wing
		type 1 basic Int
		uintLiteral 2
	member ro beak
		type 1 basic Int
		uintLiteral 238
typeSection ro ../tests/analyzer/typeSection.gPointer
	type 1 pointer {
		type 1 basic Int
	}
typeSection ro ../tests/analyzer/typeSection.hDynamicArray
	type 1 dynamicArray {
		type 1 basic Int
	}
`, test)
}
