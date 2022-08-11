package lexer

import "io"
import "github.com/sashakoshka/arf/file"
import "github.com/sashakoshka/arf/types"

// LexingOperation holds information about an ongoing lexing operataion.
type LexingOperation struct {
	file   *file.File
	char   rune
	tokens []Token
}

// Tokenize converts a file into a slice of tokens (lexemes).
func Tokenize (file *file.File) (tokens []Token, err error) {
	lexer := LexingOperation { file: file }
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
func (lexer *LexingOperation) tokenize () (err error) {
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

	if lexer.tokens[len(lexer.tokens) - 1].kind != TokenKindNewline {
		lexer.addToken(Token { kind: TokenKindNewline })
	}

	return
}

func (lexer *LexingOperation) tokenizeAlphaBeginning () (err error) {
	got := ""

	for {
		lowercase := lexer.char >= 'a' && lexer.char <= 'z'
		uppercase := lexer.char >= 'A' && lexer.char <= 'Z'
		number    := lexer.char >= '0' && lexer.char <= '9'
		if !lowercase && !uppercase && !number { break }

		got += string(lexer.char)

		lexer.nextRune()
	}

	token := Token { kind: TokenKindName, value: got }

	if len(got) == 2 {
		firstValid  := got[0] == 'n' || got[0] == 'r' || got[0] == 'w'
		secondValid := got[1] == 'n' || got[1] == 'r' || got[1] == 'w'

		if firstValid && secondValid {
			token.kind  = TokenKindPermission
			token.value = types.PermissionFrom(got)
		}
	}

	lexer.addToken(token)

	return
}

func (lexer *LexingOperation) tokenizeSymbolBeginning () (err error) {
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

		if !previousToken.Is(TokenKindNewline) ||
			!previousToken.Is(TokenKindNewline) {

			file.NewError (
				lexer.file.Location(), 1,
				"tab not used as indent",
				file.ErrorKindWarn)
			break
		}
		
		for lexer.char == '\t' {
			lexer.addToken (Token {
				kind: TokenKindIndent,
			})
			err = lexer.nextRune()
			if err != nil { return }
		}
	case '\n':
		// line break
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
		
		// TODO: if last line was blank, (ony whitespace) discard.
		lexer.addToken (Token {
			kind: TokenKindNewline,
		})
		err = lexer.nextRune()
	case '"':
		err = lexer.tokenizeString(false)
	case '\'':
		err = lexer.tokenizeString(true)
	case ':':
		lexer.addToken (Token {
			kind: TokenKindColon,
		})
		err = lexer.nextRune()
	case '.':
		lexer.addToken (Token {
			kind: TokenKindDot,
		})
		err = lexer.nextRune()
	case '[':
		lexer.addToken (Token {
			kind: TokenKindLBracket,
		})
		err = lexer.nextRune()
	case ']':
		lexer.addToken (Token {
			kind: TokenKindRBracket,
		})
		err = lexer.nextRune()
	case '{':
		lexer.addToken (Token {
			kind: TokenKindLBrace,
		})
		err = lexer.nextRune()
	case '}':
		lexer.addToken (Token {
			kind: TokenKindRBrace,
		})
		err = lexer.nextRune()
	case '+':
		err = lexer.nextRune()
		if err != nil { return }
		token := Token { kind: TokenKindPlus }
		if lexer.char == '+' {
			token.kind = TokenKindIncrement
		}
		lexer.addToken(token)
		err = lexer.nextRune()
	case '-':
		err = lexer.tokenizeDashBeginning()
	case '*':
		lexer.addToken (Token {
			kind: TokenKindAsterisk,
		})
		err = lexer.nextRune()
	case '/':
		lexer.addToken (Token {
			kind: TokenKindSlash,
		})
		err = lexer.nextRune()
	case '@':
		lexer.addToken (Token {
			kind: TokenKindAt,
		})
		err = lexer.nextRune()
	case '!':
		lexer.addToken (Token {
			kind: TokenKindExclamation,
		})
		err = lexer.nextRune()
	case '%':
		lexer.addToken (Token {
			kind: TokenKindPercent,
		})
		err = lexer.nextRune()
	case '~':
		lexer.addToken (Token {
			kind: TokenKindTilde,
		})
		err = lexer.nextRune()
	case '<':
		err = lexer.nextRune()
		if err != nil { return }
		token := Token { kind: TokenKindLessThan }
		if lexer.char == '<' {
			token.kind = TokenKindLShift
		}
		lexer.addToken(token)
		err = lexer.nextRune()
	case '>':
		err = lexer.nextRune()
		if err != nil { return }
		token := Token { kind: TokenKindGreaterThan }
		if lexer.char == '>' {
			token.kind = TokenKindRShift
		}
		lexer.addToken(token)
		err = lexer.nextRune()
	case '|':
		err = lexer.nextRune()
		if err != nil { return }
		token := Token { kind: TokenKindBinaryOr }
		if lexer.char == '|' {
			token.kind = TokenKindLogicalOr
		}
		lexer.addToken(token)
		err = lexer.nextRune()
	case '&':
		err = lexer.nextRune()
		if err != nil { return }
		token := Token { kind: TokenKindBinaryAnd }
		if lexer.char == '&' {
			token.kind = TokenKindLogicalAnd
		}
		lexer.addToken(token)
		err = lexer.nextRune()
	default:
		err = file.NewError (
			lexer.file.Location(), 1,
			"unexpected symbol character " +
			string(lexer.char),
			file.ErrorKindError)
		return
	}

	return
}

func (lexer *LexingOperation) tokenizeDashBeginning () (err error) {
	err = lexer.nextRune()
	if err != nil { return }

	if lexer.char == '-' {
		token := Token { kind: TokenKindDecrement }

		err = lexer.nextRune()
		if err != nil { return }

		if lexer.char == '-' {
			token.kind = TokenKindSeparator
			lexer.nextRune()
		}
		lexer.addToken(token)
	} else if lexer.char == '>' {
		token := Token { kind: TokenKindReturnDirection }

		err = lexer.nextRune() 
		if err != nil { return }

		lexer.addToken(token)
	} else if lexer.char >= '0' && lexer.char <= '9' {
		lexer.tokenizeNumberBeginning(true)
	} else {
		token := Token { kind: TokenKindMinus }
		lexer.addToken(token)
	}
	
	return
}

// addToken adds a new token to the lexer's token slice.
func (lexer *LexingOperation) addToken (token Token) {
	lexer.tokens = append(lexer.tokens, token)
}

// skipSpaces skips all space characters (not tabs or newlines)
func (lexer *LexingOperation) skipSpaces () (err error) {
	for lexer.char == ' ' {
		err = lexer.nextRune()
		if err != nil { return }
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
