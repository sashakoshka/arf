package parser

import "io"
import "testing"
// import "git.tebibyte.media/sashakoshka/arf/types"

func checkTree (modulePath string, correct string, test *testing.T) {
	tree, err := Parse(modulePath)
	
	if err != io.EOF {
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

// quickIdentifier returns a simple identifier of names
func quickIdentifier (trail ...string) (identifier Identifier) {
	for _, item := range trail {
		identifier.trail = append (
			identifier.trail,
			Argument {
				what:  ArgumentKindString,
				value: item,
			},
		)
	}

	return
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
