package parser

import "reflect"
import "testing"

func checkTree (modulePath string, correct *SyntaxTree, test *testing.T) {
	tree, err := Parse(modulePath)
	
	if err != nil {
		test.Log("returned error:")
		test.Log(err.Error())
		test.Fail()
		return
	}

	if !reflect.DeepEqual(tree, correct) {
		test.Log("trees not equal")
		test.Fail()
		return
	}
}

func TestMeta (test *testing.T) {
	checkTree("../tests/parser/meta",&SyntaxTree {
		license: "GPLv3",
		author:  "Sasha Koshka",
	
		requires: []string {
			"someModule",
			"otherModule",
		},
	}, test)
}
