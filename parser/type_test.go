package parser

import "testing"

func TestType (test *testing.T) {
	checkTree ("../tests/parser/type",
`:arf
---
type ro Basic:Int
type ro BasicInit:Int 6
type ro IntArray:{Int ..}
type ro IntArrayInit:{Int 3}
	3298
	923
	92
`, test)
}
