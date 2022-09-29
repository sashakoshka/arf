package analyzer

import "os"
import "testing"
import "path/filepath"
import "git.tebibyte.media/arf/arf/testCommon"

func checkTree (modulePath string, skim bool, correct string, test *testing.T) {
	cwd, _ := os.Getwd()
	modulePath = filepath.Join(cwd, modulePath)
	table, err := Analyze(modulePath, skim)
	testCommon.CheckStrings(test, table, err, correct)
}

func resolvePath (path string) (resolved string) {
	var err error
	resolved, err = filepath.Abs(path)
	if err != nil {
		panic (
			"could not resolve path (check test cases!): " +
			err.Error())
	}
	resolved = filepath.Clean(resolved)
	return
}
