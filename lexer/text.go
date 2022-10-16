package lexer

import "strconv"
import "git.tebibyte.media/arf/arf/infoerr"

// tokenizeString tokenizes a string or rune literal.
func (lexer *lexingOperation) tokenizeString () (err error) {
	token := lexer.newToken()
	
	err = lexer.nextRune()
	if err != nil { return }

	got := ""
	tokenWidth := 2

	for {
		if lexer.char == '\\' {
			err = lexer.nextRune()
			tokenWidth ++
			if err != nil { return }

			var actual     rune
			var amountRead int
			actual, amountRead, err = lexer.getEscapeSequence()
			tokenWidth += amountRead
			if err != nil { return }

			got += string(actual)
		} else {
			got += string(lexer.char)
			
			err = lexer.nextRune()
			tokenWidth ++
			if err != nil { return }
		}

		if lexer.char == '\'' { break }
	}
	
	err = lexer.nextRune()
	if err != nil { return }

	token.kind  = TokenKindString
	token.value = got

	token.location.SetWidth(tokenWidth)
	lexer.addToken(token)
	return
}

// escapeSequenceMap contains basic escape sequences and how they map to actual
// runes.
var escapeSequenceMap = map[rune] rune {
        'a':  '\x07',
        'b':  '\x08',
        'f':  '\x0c',
        'n':  '\x0a',
        'r':  '\x0d',
        't':  '\x09',
        'v':  '\x0b',
        '\'': '\'',
        '\\': '\\',
}

// getEscapeSequence reads an escape sequence in a string or rune literal.
func (lexer *lexingOperation) getEscapeSequence () (
	result     rune,
	amountRead int,
	err        error,
) {
	result, exists := escapeSequenceMap[lexer.char]
	if exists {
                err = lexer.nextRune()
                amountRead ++
                return
	} else if lexer.char >= '0' && lexer.char <= '7' {
                // octal escape sequence
                number := string(lexer.char)
        
                err = lexer.nextRune()
                amountRead ++
                if err != nil { return }
                
                for len(number) < 3 {
                        if lexer.char < '0' || lexer.char > '7' { break }

                        number += string(lexer.char)
                        
	                err = lexer.nextRune()
	                amountRead ++
	                if err != nil { return }
                }
                
                if len(number) < 3 {
			err = infoerr.NewError (
				lexer.file.Location(1),
				"octal escape sequence too short",
				infoerr.ErrorKindError)
			return
                }

                parsedNumber, _ := strconv.ParseInt(number, 8, 8)
                result = rune(parsedNumber)
                
        } else if lexer.char == 'x' || lexer.char == 'u' || lexer.char == 'U' {
                // hexidecimal escape sequence
                want := 2
                if lexer.char == 'u' { want = 4 }
                if lexer.char == 'U' { want = 8 }
        
                number := ""

                err = lexer.nextRune()
                amountRead ++
                if err != nil { return }
                
                for len(number) < want {
                        notLower := lexer.char < 'a' || lexer.char > 'f'
                        notUpper := lexer.char < 'A' || lexer.char > 'F'
                        notNum   := lexer.char < '0' || lexer.char > '9'
                        if notLower && notUpper && notNum { break }
                        
                        number += string(lexer.char)
                        
			err = lexer.nextRune()
	                amountRead ++
			if err != nil { return }
                }
                
                if len(number) < want {
			err = infoerr.NewError (
				lexer.file.Location(1),
				"hex escape sequence too short ",
				infoerr.ErrorKindError)
			return
                }

                parsedNumber, _ := strconv.ParseInt(number, 16, want * 4)
                result = rune(parsedNumber)
	} else {
		err = infoerr.NewError (
			lexer.file.Location(1),
			"unknown escape character " +
			string(lexer.char), infoerr.ErrorKindError)
		return
	}

	return
}
