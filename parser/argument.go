package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// TODO: add support for dereferences and subscripts

var validArgumentStartTokens = []lexer.TokenKind {
	lexer.TokenKindName,
	
	lexer.TokenKindInt,
	lexer.TokenKindUInt,
	lexer.TokenKindFloat,
	lexer.TokenKindString,
	
	lexer.TokenKindLBracket,
	lexer.TokenKindLParen,
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
			err = parser.nextToken()
			if err != nil { return }
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
			
			declaration := Declaration { }
			declaration.what = what
			declaration.name = identifier.trail[0]
			declaration.location = argument.Location()

			argument.kind  = ArgumentKindDeclaration
			argument.value = declaration
			
		} else {
			argument.kind  = ArgumentKindIdentifier
			argument.value = identifier
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
		
	case lexer.TokenKindLBracket:
		argument.kind  = ArgumentKindPhrase
		argument.value, err = parser.parseArgumentLevelPhrase()

	case lexer.TokenKindLParen:
		argument.kind = ArgumentKindList
		argument.value, err = parser.parseList()

	default:
		panic (
			"unimplemented argument kind " +
			parser.token.Kind().Describe())
	}

	return
}
