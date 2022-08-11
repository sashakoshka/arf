package lexer

import "github.com/sashakoshka/arf/file"

// tokenizeSymbolBeginning lexes a token that starts with a number.
func (lexer *LexingOperation) tokenizeNumberBeginning (negative bool) (err error) {
	var number uint64

	if lexer.char == '0' {
		lexer.nextRune()

		if lexer.char == 'x' {
			lexer.nextRune()
			number, err = lexer.tokenizeNumber(16)
		} else if lexer.char == 'b' {
			lexer.nextRune()
			number, err = lexer.tokenizeNumber(2)
		} else if lexer.char == '.' {
			number, err = lexer.tokenizeNumber(10)
		} else if lexer.char >= '0' && lexer.char <= '9' {
			number, err = lexer.tokenizeNumber(8)
		} else {
			return file.NewError (
				lexer.file.Location(), 1,
				"unexpected character in number literal",
				file.ErrorKindError)
		}
	} else {
		number, err = lexer.tokenizeNumber(10)
	}

	if err != nil { return }

	token := Token { }

	if negative {
		token.kind  = TokenKindInt
		token.value = int64(number) * -1
	} else {
		token.kind  = TokenKindUInt
		token.value = uint64(number)
	}
	
	lexer.addToken(token)
	return
}

func runeToDigit (char rune, radix uint64) (digit uint64, worked bool) {
	worked = true

	if char >= '0' && char <= '9' {
		digit = uint64(char - '0')
	} else if char >= 'A' && char <= 'F' {
		digit = uint64(char - 'A' + 9)
	} else if char >= 'a' && char <= 'f' {
		digit = uint64(char - 'a' + 9)
	} else {
		worked = false
	}

	if digit >= radix {
		worked = false
	}

	return
}

// tokenizeNumber reads and tokenizes a number with the specified radix.
func (lexer *LexingOperation) tokenizeNumber (radix uint64) (number uint64, err error) {
	for {
		digit, worked := runeToDigit(lexer.char, radix)
		if !worked { break }

		number *= radix
		number += digit

		err = lexer.nextRune()
		if err != nil { return }
	}
	return
}
