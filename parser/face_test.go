package parser

import "testing"

func TestFace (test *testing.T) {
	checkTree ("../tests/parser/face", false,
`:arf
---
face ro aReadWriter:Face
	read
		> into:{Byte ..}
		< read:Int
		< err:Error
	write
		> data:{Byte ..}
		< wrote:Int
		< err:Error
face ro bDestroyer:Face
	destroy
face ro cFuncInterface:Func
	> something:Int
	< someOutput:Int
	< otherOutput:String
`, test)
}
