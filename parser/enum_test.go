package parser

import "testing"

func TestEnum (test *testing.T) {
	checkTree ("../tests/parser/enum", false,
`:arf
---
enum ro AffrontToGod:Int:4
	- bird0:
		<28394 9328
		398 9>
	- bird1:
		<23 932832
		398
		2349>
	- bird2:
		<1
		2
		3
		4>
enum ro NamedColor:U32
	- red:<0xFF0000>
	- green:<0x00FF00>
	- blue:<0x0000FF>
enum ro ThisIsTerrible:Obj:
	(
	.rw x:Int
	.rw y:Int
	)
	- up:
		(
		.x:<0>
		.y:<-1>
		)
	- down:
		(
		.x:<0>
		.y:<1>)
	- left:
		(
		.x:<-1>
		.y:<0>
		)
	- right:
		(
		.x:<1>
		.y:<0>
		)
enum ro Weekday:Int
	- sunday
	- monday
	- tuesday
	- wednesday
	- thursday
	- friday
	- saturday
`, test)
}
