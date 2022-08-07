package lexer

import (
        "github.com/sashakoshka/arf/file"
)

// TokenKind is an enum represzenting what role a token has.
type TokenKind int

const (
	TokenKindNewline TokenKind = iota
	TokenKindIndent

        TokenKindSeparator
        TokenKindPermission

        TokenKindInt
        TokenKindFloat
        TokenKindString
        TokenKindRune

        TokenKindName
        TokenKindSymbol

        TokenKindColon
        TokenKindDot
        
        TokenKindLBracket
        TokenKindRBracket
        TokenKindLBrace
        TokenKindRBrace
)

// Token represents a single token. It holds its location in the file, as well
// as a value and some semantic information defining the token's role.
type Token struct {
	kind     TokenKind
	location file.Location
	value    any
}

// Kind returns the semantic role of the token.
func (token Token) Kind () (kind TokenKind) {
	return token.kind
}

// Value returns the value of the token. Depending on what kind of token it is,
// this value may be nil.
func (token Token) Value () (value any) {
	return token.value
}

// Location returns the location of the token in its file.
func (token Token) Location () (location file.Location) {
	return token.location
}
