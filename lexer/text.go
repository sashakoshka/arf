package lexer

import "github.com/sashakoshka/arf/file"

func (lexer *LexingOperation) tokenizeString (isRuneLiteral bool) (err error) {
	err = lexer.nextRune()
	if err != nil { return }

	got := ""

	for {
		got += string(lexer.char)
		
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
