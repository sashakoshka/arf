package parser

import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"

// parseData parses a data section.
func (parser *ParsingOperation) parseDataSection () (
	section *DataSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &DataSection { location: parser.token.Location() }

	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	section.what, err = parser.parseType()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindNewline) {
		err = parser.nextToken()
		if err != nil { return }

		section.value, err = parser.parseInitializationValues(0)
		if err != nil { return }
	} else {
		var argument Argument
		argument, err = parser.parseArgument()
		section.value = append(section.value, argument)

		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
	return
}

// parseInitializationValues starts on the line after a data section, or a set
// phrase. It checks for an indent greater than the indent of the aforementioned
// data section or set phrase (passed through baseIndent), and if there is,
// it parses initialization values.
func (parser *ParsingOperation) parseInitializationValues (
	baseIndent int,
) (
	values []Argument,
	err    error,
) {
	// check if line is indented one more than baseIndent
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != baseIndent + 1 { return }

	if parser.token.Is(lexer.TokenKindDot) {
		// TODO: parse as object initialization
	} else {
		// TODO: parse as array initialization
	}
	
	return
}

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
			lexer.TokenKindRBrace)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindUInt) {
			what.kind = TypeKindArray
		
			what.length = parser.token.Value().(uint64)
		
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
				file.ErrorKindError)
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
		// TODO: eat up newlines and tabs after the dot, but not before
		// it.
		if !parser.token.Is(lexer.TokenKindName) { break }

		identifier.trail = append (
			identifier.trail,
			parser.token.Value().(string))

		err = parser.nextToken()
		if err != nil { return }
		
		if !parser.token.Is(lexer.TokenKindDot) { break }
	}

	return
}
