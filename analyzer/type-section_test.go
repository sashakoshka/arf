package analyzer

import "testing"

func TestTypeSection (test *testing.T) {
	checkTree ("../tests/analyzer/typeSection", false,
`
typeSection ../tests/analyzer/typeSection.basicInt
	type basic Int
`, test)
}
