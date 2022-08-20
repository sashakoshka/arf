package parser

import "io"
import "strings"
import "testing"
// import "git.tebibyte.media/sashakoshka/arf/types"

func checkTree (modulePath string, correct string, test *testing.T) {
	tree, err  := Parse(modulePath)
	treeString := tree.ToString(0)
	treeRunes  := []rune(treeString)
	
	test.Log("CORRECT TREE:")
	logWithLineNumbers(correct, test)
	test.Log("WHAT WAS PARSED:")
	logWithLineNumbers(treeString, test)
	
	if err != io.EOF && err != nil {
		test.Log("returned error:")
		test.Log(err.Error())
		test.Fail()
		return
	}

	equal  := true
	line   := 0
	column := 0

	for index, correctChar := range correct {
		if index >= len(treeRunes) {			
			test.Log (
				"parsed is too short at line", line + 1,
				"col", column + 1)
			test.Fail()
			return
		}
		
		if correctChar != treeRunes[index] {		
			test.Log (
				"trees not equal at line", line + 1,
				"col", column + 1)
			test.Log("correct: [" + string(correctChar) + "]")
			test.Log("got:     [" + string(treeRunes[index]) + "]")
			test.Fail()
			return
		}

		if correctChar == '\n' {
			line ++
			column = 0
		} else {
			column ++
		}
	}
	
	if len(treeString) > len(correct) {
		test.Log("parsed is too long")
		test.Fail()
		return
	}

	if !equal {
		return
	}
}

func logWithLineNumbers (bigString string, test *testing.T) {
	lines := strings.Split (
		strings.Replace(bigString, "\t", "        ", -1), "\n")

	for index, line := range lines {
		test.Logf("%3d | %s", index + 1, line)
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
data ro integer:Int 3202
data ro integerArray16:{Int 16}
data ro integerArrayInitialized:{Int 16}
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
data ro integerArrayVariable:{Int ..}
data ro integerPointer:{Int}
data ro mutInteger:Int:mut 3202
data ro mutIntegerPointer:{Int}:mut
data ro nestedObject:Obj
	.that
		.bird2 123.8439
		.bird3 9328.21348239
	.this
		.bird0 324
		.bird1 "hello world"
data ro object:Obj
	.that 2139
	.this 324
`, test)
}

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
	ro this:Int
`, test)
}

