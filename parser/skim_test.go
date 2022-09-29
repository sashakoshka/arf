package parser

import "testing"

func TestSkim (test *testing.T) {
	checkTree ("../tests/parser/skim", true,
`:arf
---
data ro aExternalData:Int
	external
data ro bSingleValue:Int
	external
data ro cNestedObject:Obj
	external
data ro dUninitialized:Int:16:mut
	external
data ro eIntegerArrayInitialized:Int:16:mut
	external
func ro fComplexFunction
	---
	external
func ro gExternalFunction
	> x:Int
	< arr:Int 5
	---
	external
`, test)
}
