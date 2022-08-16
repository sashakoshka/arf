package parser

import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/lexer"

var validArgumentStartTokens = []lexer.TokenKind {
	lexer.TokenKindName,
	
	lexer.TokenKindInt,
	lexer.TokenKindUInt,
	lexer.TokenKindFloat,
	lexer.TokenKindString,
	lexer.TokenKindRune,
	
	lexer.TokenKindLBrace,
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
					file.ErrorKindError)
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
			argument.value = identifier
		}
		
	case lexer.TokenKindInt:
		
	case lexer.TokenKindUInt:
		
	case lexer.TokenKindFloat:
		
	case lexer.TokenKindString:
		
	case lexer.TokenKindRune:
		
	case lexer.TokenKindLBrace:
		
	case lexer.TokenKindLBracket:

	default:
		panic (
			"unimplemented argument kind " +
			parser.token.Kind().Describe())
	}

	return
}
