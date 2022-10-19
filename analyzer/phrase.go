package analyzer

import "git.tebibyte.media/arf/arf/parser"

type Phrase interface {
	
}

type ArbitraryPhrase struct {
	phraseBase
	command string
	arguments []Argument
}

type CastPhrase struct {
	phraseBase
	command Argument
	arguments []Argument
}

// TODO more phrases lol

func (analyzer *analysisOperation) analyzePhrase (
	inputPhrase parser.Phrase,
) (
	phrase Phrase,
	err error,
) {
	return
}
