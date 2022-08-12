package lexer

import "fmt"
import "git.tebibyte.media/sashakoshka/arf/file"

// TokenKind is an enum represzenting what role a token has.
type TokenKind int

const (
	TokenKindNewline TokenKind = iota
	TokenKindIndent

        TokenKindSeparator
        TokenKindPermission
        TokenKindReturnDirection

        TokenKindInt
        TokenKindUInt
        TokenKindFloat
        TokenKindString
        TokenKindRune

        TokenKindName

        TokenKindColon
        TokenKindDot
        
        TokenKindLBracket
        TokenKindRBracket
        TokenKindLBrace
        TokenKindRBrace

        TokenKindPlus
        TokenKindMinus
        TokenKindIncrement
        TokenKindDecrement
        TokenKindAsterisk
        TokenKindSlash

        TokenKindAt
        TokenKindExclamation
        TokenKindPercent
        TokenKindTilde

        TokenKindLessThan
        TokenKindLShift
        TokenKindGreaterThan
        TokenKindRShift
        TokenKindBinaryOr
        TokenKindLogicalOr
        TokenKindBinaryAnd
        TokenKindLogicalAnd
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

// Is returns whether or not the token is of kind kind.
func (token Token) Is (kind TokenKind) (match bool) {
	return token.kind == kind
}

// Value returns the value of the token. Depending on what kind of token it is,
// this value may be nil.
func (token Token) Value () (value any) {
	return token.value
}

// Equals returns whether this token is equal to another token
func (token Token) Equals (testToken Token) (match bool) {
	return token == testToken
} 

// Location returns the location of the token in its file.
func (token Token) Location () (location file.Location) {
	return token.location
}

// NewError creates a new error at this token's location.
func (token Token) NewError (message string, kind file.ErrorKind) (err file.Error) {
	return token.location.NewError(message, kind)
}

// Describe generates a textual description of the token to be used in debug
// logs.
func (token Token) Describe () (description string) {
	switch token.kind {
	case TokenKindNewline:
		description += "Newline"
	case TokenKindIndent:
		description += "Indent"
	case TokenKindSeparator:
		description += "Separator"
	case TokenKindPermission:
		description += "Permission"
	case TokenKindReturnDirection:
		description += "ReturnDirection"
	case TokenKindInt:
		description += "Int"
	case TokenKindUInt:
		description += "UInt"
	case TokenKindFloat:
		description += "Float"
	case TokenKindString:
		description += "String"
	case TokenKindRune:
		description += "Rune"
	case TokenKindName:
		description += "Name"
	case TokenKindColon:
		description += "Colon"
	case TokenKindDot:
		description += "Dot"
	case TokenKindLBracket:
		description += "LBracket"
	case TokenKindRBracket:
		description += "RBracket"
	case TokenKindLBrace:
		description += "LBrace"
	case TokenKindRBrace:
		description += "RBrace"
	case TokenKindPlus:
		description += "Plus"
	case TokenKindMinus:
		description += "Minus"
	case TokenKindIncrement:
		description += "Increment"
	case TokenKindDecrement:
		description += "Decrement"
	case TokenKindAsterisk:
		description += "Asterisk"
	case TokenKindSlash:
		description += "Slash"
	case TokenKindAt:
		description += "At"
	case TokenKindExclamation:
		description += "Exclamation"
	case TokenKindPercent:
		description += "Percent"
	case TokenKindTilde:
		description += "Tilde"
	case TokenKindLessThan:
		description += "LessThan"
	case TokenKindLShift:
		description += "LShift"
	case TokenKindGreaterThan:
		description += "GreaterThan"
	case TokenKindRShift:
		description += "RShift"
	case TokenKindBinaryOr:
		description += "BinaryOr"
	case TokenKindLogicalOr:
		description += "LogicalOr"
	case TokenKindBinaryAnd:
		description += "BinaryAnd"
	case TokenKindLogicalAnd:
		description += "LogicalAnd"
	}

	if token.value != nil {
		description += fmt.Sprint(": ", token.value)
	}

	return
}
