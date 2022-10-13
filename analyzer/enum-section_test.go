package analyzer

import "testing"

func TestEnumSection (test *testing.T) {
	checkTree ("../tests/analyzer/enumSection", false,
`
`, test)
}
