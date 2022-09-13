package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseType parses a type notation of the form Name, {Name}, etc.
func (parser *ParsingOperation) parseType () (what Type, err error) {
	err = parser.expect(lexer.TokenKindName, lexer.TokenKindLBrace)
	if err != nil { return }
	what.location = parser.token.Location()

	if parser.token.Is(lexer.TokenKindLBrace) {
		what.kind = TypeKindPointer

		err = parser.nextToken()
		if err != nil { return }
	
		var points Type
		points, err = parser.parseType()
		if err != nil { return }
		what.points = &points

		err = parser.expect (
			lexer.TokenKindRBrace,
			lexer.TokenKindElipsis)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindElipsis) {
			what.kind = TypeKindVariableArray
		
			err = parser.nextToken(lexer.TokenKindRBrace)
			if err != nil { return }
		}

		err = parser.nextToken()
		if err != nil { return }
	} else {
		what.name, err = parser.parseIdentifier()
		if err != nil { return }
	}

	for {
		if !parser.token.Is(lexer.TokenKindColon) { break }
		
		err = parser.nextToken(lexer.TokenKindName, lexer.TokenKindUInt)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindName) {
			// parse type qualifier
			qualifier := parser.token.Value().(string)
			switch qualifier {
			case "mut":
				what.mutable = true
			default:
				err = parser.token.NewError (
					"unknown type qualifier \"" +
					qualifier + "\"",
					infoerr.ErrorKindError)
				return
			}
		} else {
			// parse fixed array length
			what.length = parser.token.Value().(uint64)
		}
		
		err = parser.nextToken()
		if err != nil { return }
	}

	return
}
