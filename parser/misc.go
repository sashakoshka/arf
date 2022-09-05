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
			lexer.TokenKindUInt,
			lexer.TokenKindRBrace,
			lexer.TokenKindElipsis)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindUInt) {
			what.kind = TypeKindArray
		
			what.length = parser.token.Value().(uint64)
		
			err = parser.nextToken(lexer.TokenKindRBrace)
			if err != nil { return }
		} else if parser.token.Is(lexer.TokenKindElipsis) {
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

	if parser.token.Is(lexer.TokenKindColon) {
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }

		qualifier := parser.token.Value().(string)
		switch qualifier {
		case "mut":
			what.mutable = true
		default:
			err = parser.token.NewError (
				"unknown type qualifier \"" + qualifier + "\"",
				infoerr.ErrorKindError)
			return
		}
		
		err = parser.nextToken()
		if err != nil { return }
	}

	return
}

// parseIdentifier parses an identifier made out of dot separated names.
func (parser *ParsingOperation) parseIdentifier () (
	identifier Identifier,
	err        error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	identifier.location = parser.token.Location()

	for {
		if !parser.token.Is(lexer.TokenKindName) { break }

		identifier.trail = append (
			identifier.trail,
			parser.token.Value().(string))

		err = parser.nextToken()
		if err != nil { return }
		
		if !parser.token.Is(lexer.TokenKindDot) { break }

		err = parser.nextToken()
		if err != nil { return }

		// allow the identifier to continue on to the next line if there
		// is a line break right after the dot
		for parser.token.Is(lexer.TokenKindNewline) ||
			parser.token.Is(lexer.TokenKindIndent) {

			err = parser.nextToken()
			if err != nil { return }
		}
	}

	return
}
