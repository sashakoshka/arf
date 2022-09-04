package lexer

import "strconv"
import "git.tebibyte.media/arf/arf/infoerr"

// tokenizeSymbolBeginning lexes a token that starts with a number.
func (lexer *LexingOperation) tokenizeNumberBeginning (negative bool) (err error) {
	var intNumber   uint64
	var floatNumber float64
	var isFloat     bool
	var amountRead  int
	var totalRead   int

	token := lexer.newToken()

	if lexer.char == '0' {
		lexer.nextRune()
		totalRead ++

		if lexer.char == 'x' {
			lexer.nextRune()
			totalRead ++
			intNumber, floatNumber,
			isFloat, amountRead,
			err = lexer.tokenizeNumber(16)
			
		} else if lexer.char == 'b' {
			lexer.nextRune()
			totalRead ++
			intNumber, floatNumber,
			isFloat, amountRead,
			err = lexer.tokenizeNumber(2)
			
		} else if lexer.char == '.' {
			intNumber, floatNumber,
			isFloat, amountRead,
			err = lexer.tokenizeNumber(10)
			
		} else if lexer.char >= '0' && lexer.char <= '7' {
			intNumber, floatNumber,
			isFloat, amountRead,
			err = lexer.tokenizeNumber(8)
		} else if lexer.char >= '8' && lexer.char <= '9' {
			err = infoerr.NewError (
				lexer.file.Location(1),
				"unexpected number '" + string(lexer.char) +
				"' in octal literal",
				infoerr.ErrorKindError)
			return
		}
	} else {
		intNumber, floatNumber,
		isFloat, amountRead,
		err = lexer.tokenizeNumber(10)
	}

	totalRead += amountRead
	if negative {
		totalRead += 1
	}

	if err != nil { return }

	if isFloat {
		token.kind  = TokenKindFloat
		if negative {
			token.value = floatNumber * -1
		} else {
			token.value = floatNumber
		}
	} else {
		if negative {
			token.kind  = TokenKindInt
			token.value = int64(intNumber) * -1
		} else {
			token.kind  = TokenKindUInt
			token.value = uint64(intNumber)
		}		
	}

	token.location.SetWidth(totalRead)
	lexer.addToken(token)
	return
}

// runeIsDigit checks to see if the rune is a valid digit within the given
// radix, up to 16. A '.' rune will also be treated as valid.
func runeIsDigit (char rune, radix uint64) (isDigit bool) {
	isDigit = true

	var digit uint64
	if char >= '0' && char <= '9' {
		digit = uint64(char - '0')
	} else if char >= 'A' && char <= 'F' {
		digit = uint64(char - 'A' + 10)
	} else if char >= 'a' && char <= 'f' {
		digit = uint64(char - 'a' + 10)
	} else if char != '.' {
		isDigit = false
	}

	if digit >= radix {
		isDigit = false
	}

	return
}

// tokenizeNumber reads and tokenizes a number with the specified radix.
func (lexer *LexingOperation) tokenizeNumber (
	radix uint64,
) (
	intNumber   uint64,
	floatNumber float64,
	isFloat     bool,
	amountRead  int,
	err         error,
) {
	got := ""
	for {
		if !runeIsDigit(lexer.char, radix) { break }
		if lexer.char == '.' {
			if radix != 10 {
				err = infoerr.NewError (
					lexer.file.Location(1),
					"floats must have radix of 10",
					infoerr.ErrorKindError)
				return
			}
			isFloat = true
		}

		got += string(lexer.char)
		err = lexer.nextRune()
		if err != nil { return }
	}

	amountRead = len(got)

	if isFloat {
		floatNumber, err = strconv.ParseFloat(got, 64)
	} else {
		intNumber, err = strconv.ParseUint(got, int(radix), 64)
	}
	
	if err != nil {
		err = infoerr.NewError (
			lexer.file.Location(1),
			"could not parse number: " + err.Error(),
			infoerr.ErrorKindError)
		return
	}

	return
}
