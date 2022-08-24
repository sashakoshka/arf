package parser

import "testing"

func TestObjt (test *testing.T) {
	checkTree ("../tests/parser/objt",
`:arf
---
objt ro Basic:Obj
	ro that:Basic
	ro this:Basic
objt ro ComplexInit:Obj
	ro basic:Int 87
	ro complex0:Bird
		.that 98
		.this 2
	ro complex1:Bird
		.that 98902
		.this 235
	ro whatever:{Int 3}
		230984
		849
		394580
objt ro Init:Obj
	ro that:String "hello world"
	ro this:Int 23
`, test)
}
