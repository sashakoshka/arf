package analyzer

import "testing"

func TestFuncSection (test *testing.T) {
	checkTree ("../tests/analyzer/funcSection", false,
`
typeSection ro ../tests/analyzer/funcSection.aCString
	type 1 pointer {Int}
funcSection ro ../tests/analyzer/funcSection.bArbitrary
	block
		arbitraryPhrase
			command 'puts'
			cast
				type aCString
				arg string 'hellorld` + "\000" + `'
`, test)
}
