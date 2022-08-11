package lexer

import "github.com/sashakoshka/arf/file"

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

func (lexer *LexingOperation) tokenizeString (isRuneLiteral bool) (err error) {
	err = lexer.nextRune()
	if err != nil { return }

	got := ""

	for {
		// TODO: add hexadecimal escape codes
		if lexer.char == '\\' {
			err = lexer.nextRune()
			if err != nil { return }
	
			actual, exists := escapeSequenceMap[lexer.char]
			if exists {
				got += string(actual)
			} else {
				err = file.NewError (
					lexer.file.Location(), 1,
					"unknown escape character " +
					string(lexer.char), file.ErrorKindError)
				return
			}
		} else {
			got += string(lexer.char)
		}
		
		err = lexer.nextRune()
		if err != nil { return }

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
