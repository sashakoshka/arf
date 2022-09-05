package parser

import "testing"

func TestData (test *testing.T) {
	checkTree ("../tests/parser/data",
`:arf
---
data ro aInteger:Int 3202
data ro bMutInteger:Int:mut 3202
data ro cIntegerPointer:{Int}
data ro dMutIntegerPointer:{Int}:mut
data ro eIntegerArray16:Int:16
data ro fIntegerArrayVariable:{Int ..}
data ro gIntegerArrayInitialized:Int:16
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
data ro jObject:thing.Thing.thing.thing
	.this 324
	.that 2139
data ro kNestedObject:Obj
	.this
		.bird0 324
		.bird1 "hello world"
	.that
		.bird2 123.8439
		.bird3 9328.21348239
data ro lMutIntegerArray16:Int:16:mut
data ro mIntegerArrayInitialized:Int:16:mut
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
`, test)
}
