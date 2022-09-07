package parser

import "os"
import "testing"
import "path/filepath"

func TestMeta (test *testing.T) {
	cwd, _ := os.Getwd()
	checkTree ("../tests/parser/meta", false,
`:arf
author "Sasha Koshka"
license "GPLv3"
require "` + filepath.Join(cwd, "./some/local/module") + `"
require "/some/absolute/path/to/someModule"
require "/usr/include/arf/someLibraryInstalledInStandardLocation"
---
`, test)
}
