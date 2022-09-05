package parser

import "io"
import "os"
import "strings"
import "testing"
import "path/filepath"

func checkTree (modulePath string, skim bool, correct string, test *testing.T) {
	cwd, _ := os.Getwd()
	modulePath = filepath.Join(cwd, modulePath)
	tree, err := Fetch(modulePath, skim)
	
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
