package analyzer

import "git.tebibyte.media/arf/arf/parser"

// Block represents a scoped block of phrases.
type Block struct {
	phrases []Phrase

	// TODO: create a scope struct and embed it
}

func (block Block) ToString (indent int) (output string) {
	output += doIndent(indent, "block\n")

	// TODO: variables
	
	for _, phrase := range block.phrases {
		output += phrase.ToString(indent + 1)
	}
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
	for _, inputPhrase := range inputBlock {
		var outputPhrase Phrase
		outputPhrase, err = analyzer.analyzePhrase(inputPhrase)
		block.phrases = append(block.phrases, outputPhrase)
	}
	
	return
}
