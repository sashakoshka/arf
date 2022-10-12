/*
Package lexer implements a tokenizer for the ARF language. It contains a
function called Tokenize which takes in a file from the ARF file package, and
outputs an array of tokens.
*/
package lexer

import "io"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

// lexingOperation holds information about an ongoing lexing operataion.
type lexingOperation struct {
	file   *file.File
	char   rune
	tokens []Token
}

// Tokenize converts a file into a slice of tokens (lexemes).
func Tokenize (file *file.File) (tokens []Token, err error) {
	lexer := lexingOperation { file: file }
	err    = lexer.tokenize()
	tokens = lexer.tokens

	// if the lexing operation returned io.EOF, nothing went wrong so we
	// return nil for err.
	if err == io.EOF {
		err = nil
	}
	return
}

// tokenize converts a file into a slice of tokens (lexemes). It will always
// return a non-nil error, but if nothing went wrong it will return io.EOF.
func (lexer *lexingOperation) tokenize () (err error) {
	// check to see if the beginning of the file says :arf
	var shebangCheck = []rune(":arf\n")
	for index := 0; index < 5; index ++ {
		err = lexer.nextRune()
		
		if err != nil || shebangCheck[index] != lexer.char {
			err = infoerr.NewError (
				lexer.file.Location(1),
				"not an arf file",
				infoerr.ErrorKindError)
			return
		}
	}

	err = lexer.nextRune()
	if err != nil { return }

	for {
		lowercase := lexer.char >= 'a' && lexer.char <= 'z'
		uppercase := lexer.char >= 'A' && lexer.char <= 'Z'
		number    := lexer.char >= '0' && lexer.char <= '9'

		if number {
			err = lexer.tokenizeNumberBeginning(false)
			if err != nil { return }
		} else if lowercase || uppercase {
			err = lexer.tokenizeAlphaBeginning()
			if err != nil { return }
		} else {
			err = lexer.tokenizeSymbolBeginning()
			if err != nil { return }
		}

		err = lexer.skipSpaces()
		if err != nil { return }
	}

	// TODO: figure out why this is here and what its proper place is
	// because it is apparently unreachable
	if lexer.tokens[len(lexer.tokens) - 1].kind != TokenKindNewline {
		token := lexer.newToken()
		token.kind = TokenKindNewline
		lexer.addToken(token)
	}

	return
}

func (lexer *lexingOperation) tokenizeAlphaBeginning () (err error) {
	token := lexer.newToken()
	token.kind = TokenKindName

	got := ""

	for {
		lowercase := lexer.char >= 'a' && lexer.char <= 'z'
		uppercase := lexer.char >= 'A' && lexer.char <= 'Z'
		number    := lexer.char >= '0' && lexer.char <= '9'
		if !lowercase && !uppercase && !number { break }

		got += string(lexer.char)

		lexer.nextRune()
	}

	token.value = got
	token.location.SetWidth(len(got))

	if len(got) == 2 {
		permission, isPermission := types.PermissionFrom(got)
		
		if isPermission {
			token.kind  = TokenKindPermission
			token.value = permission
		}
	}

	lexer.addToken(token)

	return
}

func (lexer *lexingOperation) tokenizeSymbolBeginning () (err error) {
	switch lexer.char {
	case '#':
		// comment
		for lexer.char != '\n' {
			err = lexer.nextRune()
			if err != nil { return }
		}
	case '\t':
		// indent level
		previousToken := lexer.tokens[len(lexer.tokens) - 1]

		if !previousToken.Is(TokenKindNewline) {
			err = lexer.nextRune()
			
			infoerr.NewError (
				lexer.file.Location(1),
				"tab not used as indent",
				infoerr.ErrorKindWarn).Print()
			return
		}
		
		token := lexer.newToken()
		token.kind = TokenKindIndent

		// eat up tabs while increasing the indent level
		indentLevel := 0		
		for lexer.char == '\t' {
			indentLevel ++
			err = lexer.nextRune()
			if err != nil { return }
		}

		token.value = indentLevel
		token.location.SetWidth(indentLevel)
		lexer.addToken(token)
	case '\n':
		// line break

		// if the last line is empty, discard it
		lastLineEmpty := true
		tokenIndex := len(lexer.tokens) - 1
		for lexer.tokens[tokenIndex].kind != TokenKindNewline  {
			if lexer.tokens[tokenIndex].kind != TokenKindIndent {
				lastLineEmpty = false
				break
			}	
			tokenIndex --
		}

		if lastLineEmpty {
			lexer.tokens = lexer.tokens[:tokenIndex]
		}
		
		token := lexer.newToken()
		token.kind = TokenKindNewline
		lexer.addToken(token)
		err = lexer.nextRune()
	case '\'':
		err = lexer.tokenizeString()
	case ':':
		token := lexer.newToken()
		token.kind = TokenKindColon
		lexer.addToken(token)
		err = lexer.nextRune()
	case '.':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindDot
		if lexer.char == '.' {
			token.kind = TokenKindElipsis
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case ',':
		token := lexer.newToken()
		token.kind = TokenKindComma
		lexer.addToken(token)
		err = lexer.nextRune()
	case '(':
		token := lexer.newToken()
		token.kind = TokenKindLParen
		lexer.addToken(token)
		err = lexer.nextRune()
	case ')':
		token := lexer.newToken()
		token.kind = TokenKindRParen
		lexer.addToken(token)
		err = lexer.nextRune()
	case '[':
		token := lexer.newToken()
		token.kind = TokenKindLBracket
		lexer.addToken(token)
		err = lexer.nextRune()
	case ']':
		token := lexer.newToken()
		token.kind = TokenKindRBracket
		lexer.addToken(token)
		err = lexer.nextRune()
	case '{':
		token := lexer.newToken()
		token.kind = TokenKindLBrace
		lexer.addToken(token)
		err = lexer.nextRune()
	case '}':
		token := lexer.newToken()
		token.kind = TokenKindRBrace
		lexer.addToken(token)
		err = lexer.nextRune()
	case '+':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindPlus
		if lexer.char == '+' {
			token.kind = TokenKindIncrement
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '-':
		err = lexer.tokenizeDashBeginning()
	case '*':
		token := lexer.newToken()
		token.kind = TokenKindAsterisk
		lexer.addToken(token)
		err = lexer.nextRune()
	case '/':
		token := lexer.newToken()
		token.kind = TokenKindSlash
		lexer.addToken(token)
		err = lexer.nextRune()
	case '@':
		token := lexer.newToken()
		token.kind = TokenKindAt
		lexer.addToken(token)
		err = lexer.nextRune()
	case '!':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindExclamation
		if lexer.char == '=' {
			token.kind = TokenKindNotEqualTo
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '%':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindPercent
		if lexer.char == '=' {
			token.kind = TokenKindPercentAssignment
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '~':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindTilde
		if lexer.char == '=' {
			token.kind = TokenKindTildeAssignment
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '=':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindAssignment
		if lexer.char == '=' {
			token.kind = TokenKindEqualTo
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '<':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindLessThan
		if lexer.char == '<' {
			token.kind = TokenKindLShift
			err = lexer.nextRune()
			token.location.SetWidth(2)
			if lexer.char == '=' {
				token.kind = TokenKindLShiftAssignment
				err = lexer.nextRune()
				token.location.SetWidth(3)
			}
		} else if lexer.char == '=' {
			token.kind = TokenKindLessThanEqualTo
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '>':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindGreaterThan
		if lexer.char == '>' {
			token.kind = TokenKindRShift
			err = lexer.nextRune()
			token.location.SetWidth(2)
			if lexer.char == '=' {
				token.kind = TokenKindRShiftAssignment
				err = lexer.nextRune()
				token.location.SetWidth(3)
			}
		} else if lexer.char == '=' {
			token.kind = TokenKindGreaterThanEqualTo
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '|':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindBinaryOr
		if lexer.char == '|' {
			token.kind = TokenKindLogicalOr
			err = lexer.nextRune()
			token.location.SetWidth(2)
		} else if lexer.char == '=' {
			token.kind = TokenKindBinaryOrAssignment
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '&':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindBinaryAnd
		if lexer.char == '&' {
			token.kind = TokenKindLogicalAnd
			err = lexer.nextRune()
			token.location.SetWidth(2)
		} else if lexer.char == '=' {
			token.kind = TokenKindBinaryAndAssignment
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	case '^':
		token := lexer.newToken()
		err = lexer.nextRune()
		if err != nil { return }
		token.kind = TokenKindBinaryXor
		if lexer.char == '=' {
			token.kind = TokenKindBinaryXorAssignment
			err = lexer.nextRune()
			token.location.SetWidth(2)
		}
		lexer.addToken(token)
	default:
		err = infoerr.NewError (
			lexer.file.Location(1),
			"unexpected symbol character " +
			string(lexer.char),
			infoerr.ErrorKindError)
		return
	}

	return
}

func (lexer *lexingOperation) tokenizeDashBeginning () (err error) {
	token := lexer.newToken()
	err = lexer.nextRune()
	if err != nil { return }

	if lexer.char == '-' {
		token.kind = TokenKindDecrement
		token.location.SetWidth(2)

		err = lexer.nextRune()
		if err != nil { return }

		if lexer.char == '-' {
			token.kind = TokenKindSeparator
			lexer.nextRune()
			token.location.SetWidth(3)
		}
		lexer.addToken(token)
	} else if lexer.char == '>' {
		token.kind = TokenKindReturnDirection
		token.location.SetWidth(2)

		err = lexer.nextRune()
		if err != nil { return }

		lexer.addToken(token)
	} else if lexer.char >= '0' && lexer.char <= '9' {
		lexer.tokenizeNumberBeginning(true)
	} else {
		token.kind = TokenKindMinus
		lexer.addToken(token)
	}
	
	return
}

// newToken creates a new token from the lexer's current position in the file.
func (lexer *lexingOperation) newToken () (token Token) {
	return Token { location: lexer.file.Location(1) }
}

// addToken adds a new token to the lexer's token slice.
func (lexer *lexingOperation) addToken (token Token) {
	lexer.tokens = append(lexer.tokens, token)
}

// skipSpaces skips all space characters (not tabs or newlines)
func (lexer *lexingOperation) skipSpaces () (err error) {
	for lexer.char == ' ' {
		err = lexer.nextRune()
		if err != nil { return }
	}

	return
}

// nextRune advances the lexer to the next rune in the file.
func (lexer *lexingOperation) nextRune () (err error) {
	lexer.char, _, err = lexer.file.ReadRune()
	if err != nil && err != io.EOF {
		return infoerr.NewError (
			lexer.file.Location(1),
			err.Error(), infoerr.ErrorKindError)
	}
	return
}
