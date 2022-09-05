package parser

import "testing"

func TestMeta (test *testing.T) {
	checkTree ("../tests/parser/meta", false,
`:arf
author "Sasha Koshka"
license "GPLv3"
require "someModule"
require "otherModule"
---
`, test)
}
