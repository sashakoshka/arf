package analyzer

import "testing"

func TestEnumSection (test *testing.T) {
	checkTree ("../tests/analyzer/enumSection", false,
`enumSection ro ../tests/analyzer/enumSection.aWeekday
	arg uint 1
	type 1 basic Int
	member sunday
		arg uint 1
	member monday
		arg uint 2
	member tuesday
		arg uint 3
	member wednesday
		arg uint 4
	member thursday
		arg uint 5
	member friday
		arg uint 6
	member saturday
		arg uint 7
typeSection ro ../tests/analyzer/enumSection.bColor
	type 1 basic U32
enumSection ro ../tests/analyzer/enumSection.cNamedColor
	arg uint 16711680
	type 1 basic bColor
	member red
		arg uint 16711680
	member green
		arg uint 65280
	member blue
		arg uint 255
enumSection ro ../tests/analyzer/enumSection.dFromFarAway
	arg uint 5
	type 1 basic dInheritFromOther
	member bird
		arg uint 5
	member bread
		arg uint 4
typeSection ro ../tests/analyzer/typeSection/required.aBasic
	type 1 basic Int
typeSection ro ../tests/analyzer/typeSection.dInheritFromOther
	type 1 basic aBasic
`, test)
}
