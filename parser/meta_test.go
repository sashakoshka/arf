package parser

import "testing"

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
