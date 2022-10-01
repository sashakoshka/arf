package analyzer

import "testing"

func TestTypeSection (test *testing.T) {
	checkTree ("../tests/analyzer/typeSection", false,
`typeSection ro ../tests/analyzer/typeSection.aBasicInt
	type 1 basic Int
	arg uint 5
typeSection ro ../tests/analyzer/typeSection.bOnBasicInt
	type 1 basic aBasicInt
typeSection ro ../tests/analyzer/typeSection.cBasicObject
	type 1 basic Obj
	member ro that
		type 1 basic Int
	member ro this
		type 1 basic Int
`, test)
}
