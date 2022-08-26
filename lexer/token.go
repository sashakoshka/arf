package lexer

import "fmt"
import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

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
        TokenKindElipsis
        TokenKindComma
        
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
        TokenKindPercentAssignment
        TokenKindTilde
        TokenKindTildeAssignment

        TokenKindEqualTo
        TokenKindNotEqualTo
        TokenKindLessThanEqualTo
        TokenKindLessThan
        TokenKindLShift
        TokenKindLShiftAssignment
        TokenKindGreaterThan
        TokenKindGreaterThanEqualTo
        TokenKindRShift
        TokenKindRShiftAssignment
        TokenKindBinaryOr
        TokenKindBinaryOrAssignment
        TokenKindLogicalOr
        TokenKindBinaryAnd
        TokenKindBinaryAndAssignment
        TokenKindLogicalAnd
        TokenKindBinaryXor
        TokenKindBinaryXorAssignment
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
	return token.value == testToken.value && token.Is(testToken.kind)
} 

// Location returns the location of the token in its file.
func (token Token) Location () (location file.Location) {
	return token.location
}

// NewError creates a new error at this token's location.
func (token Token) NewError (
	message string,
	kind    infoerr.ErrorKind,
) (
	err infoerr.Error,
) {
	return infoerr.NewError(token.location, message, kind)
}

// Describe generates a textual description of the token to be used in debug
// logs.
func (token Token) Describe () (description string) {
	description = token.kind.Describe()

	if token.value != nil {
		description += fmt.Sprint(": ", token.value)
	}

	return
}

// Describe generates a textual description of the token kind to be used in
// debug logs.
func (tokenKind TokenKind) Describe () (description string) {
	switch tokenKind {
	case TokenKindNewline:
		description = "Newline"
	case TokenKindIndent:
		description = "Indent"
	case TokenKindSeparator:
		description = "Separator"
	case TokenKindPermission:
		description = "Permission"
	case TokenKindReturnDirection:
		description = "ReturnDirection"
	case TokenKindInt:
		description = "Int"
	case TokenKindUInt:
		description = "UInt"
	case TokenKindFloat:
		description = "Float"
	case TokenKindString:
		description = "String"
	case TokenKindRune:
		description = "Rune"
	case TokenKindName:
		description = "Name"
	case TokenKindColon:
		description = "Colon"
	case TokenKindDot:
		description = "Dot"
	case TokenKindElipsis:
		description = "Elipsis"
	case TokenKindComma:
		description = "Comma"
	case TokenKindLBracket:
		description = "LBracket"
	case TokenKindRBracket:
		description = "RBracket"
	case TokenKindLBrace:
		description = "LBrace"
	case TokenKindRBrace:
		description = "RBrace"
	case TokenKindPlus:
		description = "Plus"
	case TokenKindMinus:
		description = "Minus"
	case TokenKindIncrement:
		description = "Increment"
	case TokenKindDecrement:
		description = "Decrement"
	case TokenKindAsterisk:
		description = "Asterisk"
	case TokenKindSlash:
		description = "Slash"
	case TokenKindAt:
		description = "At"
	case TokenKindExclamation:
		description = "Exclamation"
	case TokenKindPercent:
		description = "Percent"
	case TokenKindPercentAssignment:
		description = "PercentAssignment"
	case TokenKindTilde:
		description = "Tilde"
	case TokenKindTildeAssignment:
		description = "TildeAssignment"
	case TokenKindEqualTo:
		description = "EqualTo"
	case TokenKindNotEqualTo:
		description = "NotEqualTo"
	case TokenKindLessThan:
		description = "LessThan"
	case TokenKindLessThanEqualTo:
		description = "LessThanEqualTo"
	case TokenKindLShift:
		description = "LShift"
	case TokenKindLShiftAssignment:
		description = "LShiftAssignment"
	case TokenKindGreaterThan:
		description = "GreaterThan"
	case TokenKindGreaterThanEqualTo:
		description = "GreaterThanEqualTo"
	case TokenKindRShift:
		description = "RShift"
	case TokenKindRShiftAssignment:
		description = "RShiftAssignment"
	case TokenKindBinaryOr:
		description = "BinaryOr"
	case TokenKindBinaryOrAssignment:
		description = "BinaryOrAssignment"
	case TokenKindLogicalOr:
		description = "LogicalOr"
	case TokenKindBinaryAnd:
		description = "BinaryAnd"
	case TokenKindBinaryAndAssignment:
		description = "BinaryAndAssignment"
	case TokenKindLogicalAnd:
		description = "LogicalAnd"
	case TokenKindBinaryXor:
		description = "BinaryXor"
	case TokenKindBinaryXorAssignment:
		description = "BinaryXorAssignment"
	}

	return
}
