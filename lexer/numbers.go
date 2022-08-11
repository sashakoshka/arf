package lexer

import "github.com/sashakoshka/arf/file"

// tokenizeSymbolBeginning lexes a token that starts with a number.
func (lexer *LexingOperation) tokenizeNumberBeginning (negative bool) (err error) {
	var number uint64

	if lexer.char == '0' {
		lexer.nextRune()

		if lexer.char == 'x' {
			lexer.nextRune()
			number, err = lexer.tokenizeHexidecimalNumber()
		} else if lexer.char == 'b' {
			lexer.nextRune()
			number, err = lexer.tokenizeBinaryNumber()
		} else if lexer.char == '.' {
			number, err = lexer.tokenizeDecimalNumber()
		} else if lexer.char >= '0' && lexer.char <= '9' {
			number, err = lexer.tokenizeOctalNumber()
		} else {
			return file.NewError (
				lexer.file.Location(), 1,
				"unexpected character in number literal",
				file.ErrorKindError)
		}
	} else {
		number, err = lexer.tokenizeDecimalNumber()
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

// tokenizeHexidecimalNumber Reads and tokenizes a hexidecimal number.
func (lexer *LexingOperation) tokenizeHexidecimalNumber () (number uint64, err error) {
	for {
		digit, worked := runeToDigit(lexer.char, 16)
		if !worked { break }

		number *= 16
		number += digit

		err = lexer.nextRune()
		if err != nil { return }
	}
	return
}

// tokenizeBinaryNumber Reads and tokenizes a binary number.
func (lexer *LexingOperation) tokenizeBinaryNumber () (number uint64, err error) {
	for {
		digit, worked := runeToDigit(lexer.char, 2)
		if !worked { break }

		number *= 2
		number += digit

		err = lexer.nextRune()
		if err != nil { return }
	}
	return
}

// tokenizeDecimalNumber Reads and tokenizes a decimal number.
func (lexer *LexingOperation) tokenizeDecimalNumber () (number uint64, err error) {
	for {
		digit, worked := runeToDigit(lexer.char, 10)
		if !worked { break }

		number *= 10
		number += digit
		
		err = lexer.nextRune()
		if err != nil { return }
	}
	
	return
}

// tokenizeOctalNumber Reads and tokenizes an octal number.
func (lexer *LexingOperation) tokenizeOctalNumber () (number uint64, err error) {
	for {
		digit, worked := runeToDigit(lexer.char, 8)
		if !worked { break }

		number *= 8
		number += digit
		
		err = lexer.nextRune()
		if err != nil { return }
	}
	return
}
