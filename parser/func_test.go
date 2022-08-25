package parser

import "testing"

func TestData (test *testing.T) {
	checkTree ("../tests/parser/func",
`:arf
---
func ro aBasicExternal
	> someInput:Int:mut
	< someOutput:Int 4
	---
	external
func ro bMethod
	@ bird:{Bird}
	> someInput:Int:mut
	< someOutput:Int 4
	---
	external
func ro cBasicPhrases
	---
	[fn 329 983 09]
	[fn 329 983 09]
	[fn 329 983 091]
	[fn [gn 329 983 091] 123]
func ro dArgumentTypes
	---
	[bird tree butterfly.wing "hello world" grass:{Int:mut 8}]
func ro eMath
	[> x:Int]
	[> y:Int]
	[< z:Int]
	[---]
	[++ x]
	[-- y]
	[set z [+ [* 0392 00] 98 x [/ 9832 y] 930]]
	[! true]
	[~ 0b01]
	[% 873 32]
	[=  5 5]
	[!= 4 4]
	[<=  4 98]
	[<   4 98]
	[<<  0x0F 4]
	[>=  98 4]
	[>   98 4]
	[>>  0xF0 4]
	[|  0b01 0b10]
	[& 0b110 0b011]
	[&& true true]
	[|| true false]
func ro fReturnDirection
	< err:Error
	---
	[someFunc 498 2980 90] -> thing:Int err
	[otherFunc] -> thing err:Error
	[fn 329 983 091] -> thing:Int err
func ro gControlFlow
	---
	[if condition]
		[something]
	[if condition]
		[something]
	[elseif]
		[otherThing]
	[else]
		[finalThing]
	[while [< x 432]]
		[something]
	[switch value]
	[: 324]
		[something]
	[: 93284]
		otherThing
	[: 9128 34738 7328]
		multipleCases
	[:]
		[defaultThing]
	[for index:Size element:Int someArray]
		[something]
		[someNextThing]
		[justMakingSureBlockParsingWorks]
	[if condition]
		[if condition]
			[nestedThing]
		[else]
			[otherThing]
	[else]
		[if condition]
			[nestedThing]
		[else]
			[otherThing]
func hSetPhrase
	---
	[set x:Int 3]
	[set y:{Int} [. x]]
	[set z:{Int 8}]
		398
		9
		2309
		983
		-2387
		478
		555
		123
	[set bird:Bird]
		.that
			.whenYou 99999
		.this 324
`, test)
}
