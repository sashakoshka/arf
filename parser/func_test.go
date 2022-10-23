package parser

import "testing"

func TestFunc (test *testing.T) {
	checkTree ("../tests/parser/func", false,
`:arf
---
func ro bMethod
	@ bird:{Bird}
	> someInput:Int:mut
	< someOutput:Int 4
	---
	external
func ro aBasicExternal
	> someInput:Int:mut
	< someOutput:Int 4
	---
	external
func ro cBasicPhrases
	---
	[fn 329 983 7]
	[fn 329 983 7]
	[fn 329 983 57]
	[fn [gn 329 983 57] 123]
func ro dArgumentTypes
	---
	[bird tree butterfly.wing 'hello world' grass:Int:8:mut]
func ro eMath
	> x:Int
	> y:Int
	< z:Int
	---
	[++ x]
	[-- y]
	[= z [+ [* 250 0] 98 x [/ 9832 y] 930]]
	[== 4 4]
	[! true]
	[~ 1]
	[~= x]
	[% 873 32]
	[= 5 5]
	[!= 4 4]
	[<= 4 98]
	[< 4 98]
	[<< 15 4]
	[<<= x 4]
	[>= 98 4]
	[> 98 4]
	[>> 240 4]
	[>>= x 4]
	[| 1 2]
	[|= x 2]
	[& 6 3]
	[&= x 3]
	[&& true true]
	[|| true false]
func ro fReturnDirection
	< err:Error
	---
	[someFunc 498 2980 90] -> thing:Int err
	[otherFunc] -> thing err:Error
	[fn 329 983 57] -> thing:Int err
func ro gControlFlow
	---
	[defer]
		[something]
		[otherThing]
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
		[otherThing]
	[: 9128 34738 7328]
		[multipleCases]
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
func ro hDataInit
	---
	[= x:Int 3]
	[= y:{Int} [loc x]]
	[= z:Int:8 (398 9 2309 983 -2387 478 555 123)]
	[= bird:Bird ((99999) 324)]
func ro iDereference
	> x:{Int}
	> y:{Int ..}
	> z:Int:4
	---
	[= b:Int {x}]
	[= c:Int {y 4}]
	[= d:Int {z 3}]
`, test)
}
