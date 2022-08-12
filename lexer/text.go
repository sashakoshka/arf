package lexer

import "strconv"
import "git.tebibyte.media/sashakoshka/arf/file"

// tokenizeString tokenizes a string or rune literal.
func (lexer *LexingOperation) tokenizeString (isRuneLiteral bool) (err error) {
	err = lexer.nextRune()
	if err != nil { return }

	got := ""

	for {
		// TODO: add hexadecimal escape codes
		if lexer.char == '\\' {
			err = lexer.nextRune()
			if err != nil { return }

			var actual rune
			actual, err = lexer.getEscapeSequence()
			if err != nil { return }

			got += string(actual)
		} else {
			got += string(lexer.char)
			
			err = lexer.nextRune()
			if err != nil { return }
		}

		if isRuneLiteral {
			if lexer.char == '\'' { break }
		} else {
			if lexer.char == '"'  { break }
		}
	}
	
	err = lexer.nextRune()
	if err != nil { return }

	token := Token { }

	if isRuneLiteral {
		if len(got) > 1 {
			err = file.NewError (
				lexer.file.Location(), len(got) - 1,
				"excess data in rune literal",
				file.ErrorKindError)
			return
		}

		token.kind  = TokenKindRune
		token.value = rune([]rune(got)[0])
	} else {
		token.kind  = TokenKindString
		token.value = got
	}

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
        '"':  '"',
        '\\': '\\',
}

// getEscapeSequence reads an escape sequence in a string or rune literal.
func (lexer *LexingOperation) getEscapeSequence () (result rune, err error) {
	result, exists := escapeSequenceMap[lexer.char]
	if exists {
                err = lexer.nextRune()
                return
	} else if lexer.char >= '0' && lexer.char <= '7' {
                // octal escape sequence
                number := string(lexer.char)
        
                err = lexer.nextRune()
                if err != nil { return }
                
                for len(number) < 3 {
                        if lexer.char < '0' || lexer.char > '7' { break }

                        number += string(lexer.char)
                        
	                err = lexer.nextRune()
	                if err != nil { return }
                }
                
                if len(number) < 3 {
			err = file.NewError (
				lexer.file.Location(), 1,
				"octal escape sequence too short",
				file.ErrorKindError)
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
                if err != nil { return }
                
                for len(number) < want {
                        notLower := lexer.char < 'a' || lexer.char > 'f'
                        notUpper := lexer.char < 'A' || lexer.char > 'F'
                        notNum   := lexer.char < '0' || lexer.char > '9'
                        if notLower && notUpper && notNum { break }
                        
                        number += string(lexer.char)
                        
			err = lexer.nextRune()
			if err != nil { return }
                }
                
                if len(number) < want {
			err = file.NewError (
				lexer.file.Location(), 1,
				"hex escape sequence too short ",
				file.ErrorKindError)
			return
                }

                parsedNumber, _ := strconv.ParseInt(number, 16, want * 4)
                result = rune(parsedNumber)
	} else {
		err = file.NewError (
			lexer.file.Location(), 1,
			"unknown escape character " +
			string(lexer.char), file.ErrorKindError)
		return
	}

	return
}
