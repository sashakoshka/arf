package analyzer

import "testing"

func TestDataSection (test *testing.T) {
	checkTree ("../tests/analyzer/dataSection", false,
`
`, test)
}
