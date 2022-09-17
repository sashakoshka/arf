package testCommon

import "io"
import "strings"
import "testing"

// Stringable encompasses any type that can be converted into a string, given
// an indentation level.
type Stringable interface {
	ToString (indent int) (output string)
}

// CheckStrings takes in an error, a result value that is able to provide a
// string, and a correct string to check the result value's string against, and
// a test pointer. This function will detect the presence of an error, or
// differences in the result string and the correct string, and fail the test
// if needed. It also provides detailed output logs to facilitate the correction
// of test failues.
func CheckStrings (
	test    *testing.T,
	result  Stringable,
	err     error,
	correct string,
) {
	treeString := result.ToString(0)
	treeRunes  := []rune(treeString)
	
	test.Log("CORRECT:")
	logWithLineNumbers(correct, test)
	test.Log("RESULT:")
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
				"result is too short at line", line + 1,
				"col", column + 1)
			test.Fail()
			return
		}
		
		if correctChar != treeRunes[index] {		
			test.Log (
				"strings not equal at line", line + 1,
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
		test.Log("result is too long")
		test.Fail()
		return
	}

	if !equal {
		return
	}
}

// logWithLineNumbers logs a multi-line, indented string with line numbers. It
// also converts all tabs to spaces for consistent indentation.
func logWithLineNumbers (bigString string, test *testing.T) {
	lines := strings.Split (
		strings.Replace(bigString, "\t", "        ", -1), "\n")

	for index, line := range lines {
		test.Logf("%3d | %s", index + 1, line)
	}
}
