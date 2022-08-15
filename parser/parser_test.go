package parser

import "io"
import "testing"
// import "git.tebibyte.media/sashakoshka/arf/types"

func checkTree (modulePath string, correct string, test *testing.T) {
	tree, err := Parse(modulePath)
	
	if err != io.EOF && err != nil {
		test.Log("returned error:")
		test.Log(err.Error())
		test.Fail()
		return
	}

	treeString := tree.ToString(0)
	if treeString != correct {
		test.Log("trees not equal!")
		test.Log("CORRECT TREE:")
		test.Log(correct)
		test.Log("WHAT WAS PARSED:")
		test.Log(treeString)
		test.Fail()
		return
	}
}

func TestMeta (test *testing.T) {
	checkTree ("../tests/parser/meta",
`:arf
author "Sasha Koshka"
license "GPLv3"
require "someModule"
require "otherModule"
---
`, test)
}

func TestData (test *testing.T) {
	checkTree ("../tests/parser/data",
`:arf
---
data wr integer:Int 3202
data wr integerPointer:{Int}
data wr integerArray16:{Int 16}
data wr integerArrayVariable:{Int ...}
data wr integerArrayInitialized:{Int 16}
	3948
	293
	293049
	948
	912
	340
	0
	2304
	0
	4785
	92
data wr object:Obj
	, this 324
	, that 2139
data wr nestedObject:Obj
	, this
		, bird0 324
		, bird1 "hello world"
	, that
		, bird2 123.8439
		, bird3 9328.21348239
`, test)
}

