package parser

import "testing"

func TestType (test *testing.T) {
	checkTree ("../tests/parser/type", false,
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
type ro eInitAndDefine:aBasic
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

type ro fBasic:Int
type ro gBasicInit:Int 6
type ro hIntArray:{Int ..}
type ro iIntArrayInit:Int:3
	3298
	923
	92
`, test)
}
