package parser

import "testing"

// TODO: merge this test with the type test
func TestObjt (test *testing.T) {
	checkTree ("../tests/parser/objt", false,
`:arf
---
type ro aBasic:Obj
	ro that:Int
	ro this:Int
type ro bBitFields:Obj
	ro that:Int & 1
	ro this:Int & 24 298
type ro cInit:Obj
	ro that:String "hello world"
	ro this:Int 23
type ro dInitInherit:aBasic
	-- that 9384
	-- this 389
type ro cInitAndDefine:aBasic
	-- this 389
	ro these:aBasic
		ro born:Int 4
		ro in:Int
		ro the:Int:3
			9348
			92384
			92834
		-- this 98
	-- that 9384
`, test)
}
