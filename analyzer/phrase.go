package analyzer

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

type Phrase interface {
	// Provided by phraseBase
	Location () (location file.Location)
	NewError (message string, kind infoerr.ErrorKind) (err error)

	// Must be implemented by each individual phrase
	ToString (indent int) (output string)
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
