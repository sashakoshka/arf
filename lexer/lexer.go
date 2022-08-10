package lexer

import "io"
import "github.com/sashakoshka/arf/file"

// LexingOperation holds information about an ongoing lexing operataion.
type LexingOperation struct {
	file   *file.File
	char   rune
	tokens []Token
}

// Tokenize converts a file into a slice of tokens (lexemes).
func Tokenize (file *file.File) (tokens []Token, err error) {
	lexer := LexingOperation { file: file }
	err = lexer.tokenize()
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
			// TODO: tokenize number begin
		} else if lowercase || uppercase {
			// TODO: tokenize alpha begin
		} else {
			err = lexer.tokenizeSymbolBeginning()
			if err != nil { return err }
		}

		// TODO:  skip whitespace
	}

	return
}

func (lexer *LexingOperation) tokenizeSymbolBeginning () (err error) {
	switch lexer.char {
	case '#':
		// comment
		for lexer.char != '\n' {
			lexer.nextRune()
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
			lexer.nextRune()
		}
	case '\n':
		// line break
		lexer.addToken (Token {
			kind: TokenKindNewline,
		})
		lexer.nextRune()
	case '"':
		// TODO: tokenize string literal
		lexer.nextRune()
	case '\'':
		// TODO: tokenize rune literal
		lexer.nextRune()
	case ':':
		lexer.addToken (Token {
			kind: TokenKindColon,
		})
		lexer.nextRune()
	case '.':
		lexer.addToken (Token {
			kind: TokenKindDot,
		})
		lexer.nextRune()
	case '[':
		lexer.addToken (Token {
			kind: TokenKindLBracket,
		})
		lexer.nextRune()
	case ']':
		lexer.addToken (Token {
			kind: TokenKindRBracket,
		})
		lexer.nextRune()
	case '{':
		lexer.addToken (Token {
			kind: TokenKindLBrace,
		})
		lexer.nextRune()
	case '}':
		lexer.addToken (Token {
			kind: TokenKindRBrace,
		})
		lexer.nextRune()
	case '+':
		lexer.addToken (Token {
			kind: TokenKindPlus,
		})
		lexer.nextRune()
	case '-':
		// TODO: tokenize dash begin
		lexer.nextRune()
	case '*':
		lexer.addToken (Token {
			kind: TokenKindAsterisk,
		})
		lexer.nextRune()
	case '/':
		lexer.addToken (Token {
			kind: TokenKindSlash,
		})
		lexer.nextRune()
	case '@':
		lexer.addToken (Token {
			kind: TokenKindAt,
		})
		lexer.nextRune()
	case '!':
		lexer.addToken (Token {
			kind: TokenKindExclamation,
		})
		lexer.nextRune()
	case '%':
		lexer.addToken (Token {
			kind: TokenKindPercent,
		})
		lexer.nextRune()
	case '~':
		lexer.addToken (Token {
			kind: TokenKindTilde,
		})
		lexer.nextRune()
	case '<':
		// TODO: tokenize less than begin
		lexer.nextRune()
	case '>':
		// TODO: tokenize greater than begin
		lexer.nextRune()
	case '|':
		// TODO: tokenize bar begin
		lexer.nextRune()
	case '&':
		// TODO: tokenize and begin
		lexer.nextRune()
	default:
		err = file.NewError (
			lexer.file.Location(), 1,
			"unexpected character " +
			string(lexer.char),
			file.ErrorKindError)
		return
	}

	return
}

func (lexer *LexingOperation) addToken (token Token) {
	lexer.tokens = append(lexer.tokens, token)
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
