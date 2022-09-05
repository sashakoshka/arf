package parser

import "testing"

func TestEnum (test *testing.T) {
	checkTree ("../tests/parser/enum", false,
`:arf
---
enum ro AffrontToGod:Int:4
	bird0
		28394
		9328
		398
		9
	bird1
		23
		932832
		398
		2349
	bird2
		1
		2
		3
		4
enum ro NamedColor:U32
	red 16711680
	green 65280
	blue 255
enum ro Weekday:Int
	sunday
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
`, test)
}
