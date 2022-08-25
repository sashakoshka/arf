package parser

import "testing"

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
data ro object:thing.thing.thing.thing
	.that 2139
	.this 324
`, test)
}
