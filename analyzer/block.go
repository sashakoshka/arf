package analyzer

import "git.tebibyte.media/arf/arf/parser"

// Block represents a scoped block of phrases.
type Block struct {
	locatable
	phrases []Phrase

	// TODO: create a scope struct and embed it
}

func (block Block) ToString (indent int) (output string) {
	output += doIndent(indent, "block\n")

	// TODO: variables
	// TODO: phrases
	return
}

// analyzeBlock analyzes a scoped block of phrases.
// TODO: have a way to "start out" with a list of variables for things like
// arguments, and declarations inside of control flow statements
func (analyzer *analysisOperation) analyzeBlock (
	inputBlock parser.Block,
) (
	block Block,
	err error,
) {
	return
}
