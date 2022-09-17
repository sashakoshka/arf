package parser

import "os"
import "testing"
import "path/filepath"
import "git.tebibyte.media/arf/arf/testCommon"

func checkTree (modulePath string, skim bool, correct string, test *testing.T) {
	cwd, _ := os.Getwd()
	modulePath = filepath.Join(cwd, modulePath)
	tree, err := Fetch(modulePath, skim)
	testCommon.CheckStrings(test, tree, err, correct)
}
