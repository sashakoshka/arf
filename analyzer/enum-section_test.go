package analyzer

import "testing"

func TestEnumSection (test *testing.T) {
	checkTree ("../tests/analyzer/enumSection", false,
`enumSection ro ../tests/analyzer/enumSection.aWeekday
	uintLiteral 1
	type 1 basic Int
	member sunday
		uintLiteral 1
	member monday
		uintLiteral 2
	member tuesday
		uintLiteral 3
	member wednesday
		uintLiteral 4
	member thursday
		uintLiteral 5
	member friday
		uintLiteral 6
	member saturday
		uintLiteral 7
typeSection ro ../tests/analyzer/enumSection.bColor
	type 1 basic U32
enumSection ro ../tests/analyzer/enumSection.cNamedColor
	uintLiteral 16711680
	type 1 basic bColor
	member red
		uintLiteral 16711680
	member green
		uintLiteral 65280
	member blue
		uintLiteral 255
enumSection ro ../tests/analyzer/enumSection.dFromFarAway
	uintLiteral 5
	type 1 basic dInheritFromOther
	member bird
		uintLiteral 5
	member bread
		uintLiteral 4
typeSection ro ../tests/analyzer/typeSection/required.aBasic
	type 1 basic Int
typeSection ro ../tests/analyzer/typeSection.dInheritFromOther
	type 1 basic aBasic
`, test)
}
