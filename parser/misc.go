package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// TODO: need to come up with another syntax for bitfields, and use that syntax
// for fixed-length arrays. the problem is, we cant use {} for those kinds of
// arrays because they aren't pointers under the hood and that goes against the
// arf design goals in a nasty ugly way, and not only that. but the :mut type
// qualifier is meaningless on fixed length arrays now. the bit field syntax
// solves both of these problems very gracefully. but now, the problem is coming
// up with new bit field syntax. implementing this change is extremely
// necessary, for it will supercharge the coherency of the language and make it
// way more awesome.
//
// this new syntax cannot conflict with arguments, so it cannot have any
// tokens that begin those. basically all symbol tokens are on the table here.
// some ideas:
//
// ro member:Type ~ 4
// ro member:Type & 4    <- i like this one because binary &. so intuitive.
// ro member:Type % 4
// ro member:Type | 4
// ro member:Type ! 4
// ro member:Type (4)    <- this looks phenomenal, but it needs new tokens not
//                          used anywhere else, and it would be mildly annoying
//                          to parse.

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
