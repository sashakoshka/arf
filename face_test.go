package parser

import "testing"

func TestFace (test *testing.T) {
	checkTree ("../tests/parser/face",
`:arf
---
face ro Destroyer:Face
	destroy
face ro ReadWriter:Face
	read
		> into:{Byte ..}
		< read:Int
		< err:Error
	write
		> data:{Byte ..}
		< wrote:Int
		< err:Error
`, test)
}
