package analyzer

import "regexp"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

var validNameRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

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

func (phrase ArbitraryPhrase) ToString (indent int) (output string) {
	output += doIndent(indent, "phrase\n")
	output += doIndent(indent + 1, phrase.command, "\n")

	for _, argument := range phrase.arguments {
		output += argument.ToString(indent + 1)
	}

	return
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
	base := phraseBase { }
	base.location = inputPhrase.Location()

	arguments := []Argument { }
	for index := 0; index < inputPhrase.Length(); index ++ {
		inputArgument := inputPhrase.Argument(index)
		
		var argument Argument
		argument, err = analyzer.analyzeArgument(inputArgument)
		if err != nil { return }
		
		arguments = append(arguments, argument)
	}

	switch inputPhrase.Kind() {
	case parser.PhraseKindArbitrary:
		command := inputPhrase.Command().Value().(string)
		if !validNameRegex.Match([]byte(command)) {
			err = inputPhrase.NewError (
				"command cannot contain characters other " +
				"than a-z, A-Z, 0-9, underscores, or begin " +
				"with a number",
				infoerr.ErrorKindError)
			return
		}

		outputPhrase := ArbitraryPhrase {
			phraseBase: base,
			command:    command,
			arguments:  arguments,
		}
		phrase = outputPhrase
		
	default:
		panic("phrase kind not implemented")
	}
	return
}
