package parser

import "testing"

func TestFace (test *testing.T) {
	checkTree ("../tests/parser/face", false,
`:arf
---
face ro ReadWriter:Face
	write
		> data:{Byte ..}
		< wrote:Int
		< err:Error
	read
		> into:{Byte ..}
		< read:Int
		< err:Error
face ro Destroyer:Face
	destroy
face ro cFuncInterface
	> something:Int
	< someOutput:Int
	< otherOutput:String
`, test)
}
