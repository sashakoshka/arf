package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

var validArgumentStartTokens = []lexer.TokenKind {
	lexer.TokenKindName,
	
	lexer.TokenKindInt,
	lexer.TokenKindUInt,
	lexer.TokenKindFloat,
	lexer.TokenKindString,
	lexer.TokenKindRune,
	
	lexer.TokenKindLBracket,
}

func (parser *ParsingOperation) parseArgument () (argument Argument, err error) {
	argument.location = parser.token.Location()

	err = parser.expect(validArgumentStartTokens...)
	if err != nil { return }

	switch parser.token.Kind() {
	case lexer.TokenKindName:
		var identifier Identifier
		identifier, err = parser.parseIdentifier()
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindColon) {
			var what Type
			what, err = parser.parseType()
			if err != nil { return }

			if len(identifier.trail) != 1 {
				err = parser.token.NewError (
					"cannot use member selection in " +
					"a variable definition",
					infoerr.ErrorKindError)
				return
			}

			argument.kind  = ArgumentKindDeclaration
			argument.value = Declaration {
				location: argument.location,
				name:     identifier.trail[0],
				what:     what,
			}
		} else {
			argument.kind  = ArgumentKindIdentifier
			argument.value = &identifier
		}
		
	case lexer.TokenKindInt:
		argument.kind  = ArgumentKindInt
		argument.value = parser.token.Value().(int64)
		err = parser.nextToken()
		
	case lexer.TokenKindUInt:
		argument.kind  = ArgumentKindUInt
		argument.value = parser.token.Value().(uint64)
		err = parser.nextToken()
		
	case lexer.TokenKindFloat:
		argument.kind  = ArgumentKindFloat
		argument.value = parser.token.Value().(float64)
		err = parser.nextToken()
		
	case lexer.TokenKindString:
		argument.kind  = ArgumentKindString
		argument.value = parser.token.Value().(string)
		parser.nextToken()
		
	case lexer.TokenKindRune:
		argument.kind  = ArgumentKindRune
		argument.value = parser.token.Value().(rune)
		parser.nextToken()
		
	case lexer.TokenKindLBracket:
		argument.kind = ArgumentKindPhrase
		var phrase Phrase
		phrase, err = parser.parseArgumentLevelPhrase()
		argument.value = &phrase
		parser.nextToken()
		

	default:
		panic (
			"unimplemented argument kind " +
			parser.token.Kind().Describe())
	}

	return
}
