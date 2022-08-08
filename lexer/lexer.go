package lexer

import "github.com/sashakoshka/arf/file"

// LexingOperation holds information about an ongoing lexing operataion.
type LexingOperation struct {
	file *file.File
}

// Tokenize converts a file into a slice of tokens (lexemes)
func Tokenize (file *file.File) (tokens []Token) {
	lexer := LexingOperation { }
	return lexer.tokenize(file)
}

// tokenize converts a file into a slice of tokens (lexemes)
func (lexer *LexingOperation) tokenize (file *file.File) (tokens []Token) {
	return
}
