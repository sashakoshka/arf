package analyzer

import "testing"

func TestDataSection (test *testing.T) {
	checkTree ("../tests/analyzer/dataSection", false,
`dataSection ro ../tests/analyzer/dataSection.aBasicInt
	type 1 basic Int
	arg uint 5
dataSection ro ../tests/analyzer/dataSection.bRune
	type 1 basic Int
	arg string 'A'
dataSection ro ../tests/analyzer/dataSection.cString
	type 1 basic String
	arg string 'A very large bird'
dataSection ro ../tests/analyzer/dataSection.dCharBuffer
	type 32 basic U8
	arg string 'A very large bird` + "\000" + `'
`, test)
}
