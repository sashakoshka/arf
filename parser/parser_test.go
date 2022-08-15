package parser

import "reflect"
import "testing"
import "git.tebibyte.media/sashakoshka/arf/types"

func checkTree (modulePath string, correct *SyntaxTree, test *testing.T) {
	tree, err := Parse(modulePath)
	
	if err != nil {
		test.Log("returned error:")
		test.Log(err.Error())
		test.Fail()
		return
	}

	if !reflect.DeepEqual(tree, correct) {
		test.Log("trees not equal")
		test.Fail()
		return
	}
}

// quickIdentifier returns a simple identifier of names
func quickIdentifier (trail ...string) (identifier Identifier) {
	for _, item := range trail {
		identifier.trail = append (
			identifier.trail,
			Argument {
				what:  ArgumentKindString,
				value: item,
			},
		)
	}

	return
}

func TestMeta (test *testing.T) {
	checkTree ("../tests/parser/meta", &SyntaxTree {
		license: "GPLv3",
		author:  "Sasha Koshka",
	
		requires: []string {
			"someModule",
			"otherModule",
		},
	}, test)
}

func TestData (test *testing.T) {
	tree := &SyntaxTree {
		dataSections: []DataSection {
			DataSection {
				name: "integer",
				permission: types.PermissionFrom("wr"),
				
				what: Type {
					kind: TypeKindBasic,
					name: quickIdentifier("Int"),
				},
				value: []Argument {
					Argument {
						what: ArgumentKindUInt,
						value: 3202,
					},
				},
			},

			DataSection {
				name: "integerPointer",
				permission: types.PermissionFrom("wr"),
				
				what: Type {
					kind: TypeKindPointer,
					points: &Type {
						kind: TypeKindBasic,
						name: quickIdentifier("Int"),
					},
				},
				value: []Argument {
					Argument {
						what:  ArgumentKindUInt,
						value: 3202,
					},
				},
			},

			DataSection {
				name: "integerArray16",
				permission: types.PermissionFrom("wr"),
				
				what: Type {
					kind: TypeKindArray,
					points: &Type {
						kind: TypeKindBasic,
						name: quickIdentifier("Int"),
					},
					length: Argument {
						what:  ArgumentKindUInt,
						value: 16,
					}
				},
				value: []Argument {
					Argument {
						what:  ArgumentKindUInt,
						value: 3202,
					},
				},
			},
		},
	}
	
	checkTree ("../tests/parser/data", tree, test)
}
