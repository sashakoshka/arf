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
		section.value, err = parser.parseArgument()
		if err != nil { return }

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
	initializationArgument Argument,
	err error,
) {
	// check if line is indented one more than baseIndent
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != baseIndent + 1 { return }
	
	err = parser.nextToken()
	if err != nil { return }

	initializationArgument.location = parser.token.Location()

	if parser.token.Is(lexer.TokenKindDot) {
		var initializationValues ObjectInitializationValues
		initializationValues, err = parser.parseObjectInitializationValues (
			baseIndent + 1)
		initializationArgument.kind = ArgumentKindObjectInitializationValues
		initializationArgument.value = &initializationValues
	} else {
		var initializationValues ArrayInitializationValues
		initializationValues, err = parser.parseArrayInitializationValues (
			baseIndent + 1)
		initializationArgument.kind = ArgumentKindArrayInitializationValues
		initializationArgument.value = &initializationValues
	}
	
	return
}

// parseObjectInitializationValues parses a list of object initialization
// values until the indentation level is not equal to indent.
func (parser *ParsingOperation) parseObjectInitializationValues (
	indent int,
) (
	values ObjectInitializationValues,
	err    error,
) {
	values.attributes = make(map[string] Argument)
	values.location   = parser.token.Location()

	for {
		// get attribute name and value
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		name := parser.token.Value().(string)
		
		_, exists := values.attributes[name]
		if exists {
			err = parser.token.NewError (
				"duplicate member \"" + name + "\" in object " +
				"member initialization",
				file.ErrorKindError)
			return
		}
		
		err = parser.nextToken()
		if err != nil { return }
		var value Argument
		value, err = parser.parseArgument()

		// store in object
		values.attributes[name] = value
		
		// go onto the next line
		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
		if !parser.token.Is(lexer.TokenKindIndent) { break }
		// TODO: if indent is greater, recurse instead
		if parser.token.Value().(int) != indent { break }

		// the next line must start with a dot
		err = parser.nextToken(lexer.TokenKindDot)
		if err != nil { return }
	}
	return
}

// parseArrayInitializationValues parses a list of array initialization values
// until the indentation lexel is not equal to indent.
func (parser *ParsingOperation) parseArrayInitializationValues (
	indent int,
) (
	values ArrayInitializationValues,
	err    error,
) {
	values.location = parser.token.Location()
	
	for {
		if parser.token.Is(lexer.TokenKindNewline) {
			err = parser.nextToken()
			if err != nil { return }
			
			if !parser.token.Is(lexer.TokenKindIndent) { break }
			if parser.token.Value().(int) != indent { break }
			err = parser.nextToken()
			if err != nil { return }
		}
		
		var argument Argument
		argument, err = parser.parseArgument()
		if err != nil { return }
		values.values = append(values.values, argument)
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
			lexer.TokenKindRBrace,
			lexer.TokenKindElipsis)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindUInt) {
			what.kind = TypeKindArray
		
			what.length = parser.token.Value().(uint64)
		
			err = parser.nextToken(lexer.TokenKindRBrace)
			if err != nil { return }
		} else if parser.token.Is(lexer.TokenKindElipsis) {
			what.kind = TypeKindArray
		
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
