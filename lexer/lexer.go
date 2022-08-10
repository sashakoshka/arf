package lexer

import "io"
import "github.com/sashakoshka/arf/file"

// LexingOperation holds information about an ongoing lexing operataion.
type LexingOperation struct {
	file *file.File
	char rune
}

// Tokenize converts a file into a slice of tokens (lexemes).
func Tokenize (file *file.File) (tokens []Token, err error) {
	lexer := LexingOperation { file: file }
	tokens, err = lexer.tokenize()

	// if the lexing operation returned io.EOF, nothing went wrong so we
	// return nil for err.
	if err == io.EOF {
		err = nil
	}
	return
}

// tokenize converts a file into a slice of tokens (lexemes). It will always
// return a non-nil error, but if nothing went wrong it will return io.EOF.
func (lexer *LexingOperation) tokenize () (tokens []Token, err error) {
	err = lexer.nextRune()
	if err != nil { return }

	for {
		lowercase := lexer.char >= 'a' && lexer.char <= 'z'
		uppercase := lexer.char >= 'A' && lexer.char <= 'Z'
		number    := lexer.char >= '0' && lexer.char <= '9'

		if number {
			// TODO: tokenize number
		} else if lowercase || uppercase {
			// TODO: tokenize multi
		} else {
			switch lexer.char {
			case '"':
				// TODO: tokenize string literal
				lexer.nextRune()
			case '\'':
				// TODO: tokenize rune literal
				lexer.nextRune()
			case ':':
				// TODO: colon token
			case '.':
				// TODO: dot token
			case '[':
				// TODO: left bracket token
			case ']':
				// TODO: right bracket token
			case '{':
				// TODO: left brace token
			case '}':
				// TODO: right brace token
			// TODO: add more for things like math symbols, return
			// direction operators, indentation, etc
			default:
				err = file.NewError (
					lexer.file.Location(), 1,
					"unexpected character " +
					string(lexer.char),
					file.ErrorKindError)
				return
			}
		}

		// TODO: skip whitespace
	}

	return
}

// nextRune advances the lexer to the next rune in the file.
func (lexer *LexingOperation) nextRune () (err error) {
	lexer.char, _, err = lexer.file.ReadRune()
	if err != nil && err != io.EOF {
		return file.NewError (
			lexer.file.Location(), 1,
			err.Error(), file.ErrorKindError)
	}
	return
}

// 
