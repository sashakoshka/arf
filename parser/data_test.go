package parser

import "testing"

func TestData (test *testing.T) {
	checkTree ("../tests/parser/data", false,
`:arf
---
data ro aInteger:Int:<3202>
data ro bMutInteger:Int:mut:<3202>
data ro cIntegerPointer:{Int}
data ro dMutIntegerPointer:{Int}:mut
data ro eIntegerArray16:Int:16
data ro fIntegerArrayVariable:{Int ..}
data ro gIntegerArrayInitialized:Int:16:
	<
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
	>
data wr hIntegerPointerInit:{Int}:<[& integer]>
data wr iMutIntegerPointerInit:{Int}:mut:<[& integer]>
data ro jObject:Obj:
	(
	.that:<324>
	.this:<324>
	)
data ro kNestedObject:Obj:
	(
	.ro newMember:Int:<9023>
	):
	(
	.that:
		(
		.bird2:<123.8439>
		.bird3:<9328.21348239>
		)
	.this:
		(
		.bird0:<324>
		.bird1:<"hello world">
		)
	)
data ro lMutIntegerArray16:Int:16:mut
data ro mExternalData:Int:8
	external
data ro nIntegerArrayInitialized:Int:16:mut:
	<
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
	>
`, test)
}
