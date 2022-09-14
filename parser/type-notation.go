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

	// TODO: consider offloading array default values to argument parsing, and
	// then just grabbing the next argument after this if it exists. that way
	// we can have array litreals anywhere.

	// get default value
	if parser.token.Is(lexer.TokenKindLessThan) {
		what.defaultValue, err = parser.parseBasicDefaultValue()
		if err != nil { return }
		
	} else if parser.token.Is(lexer.TokenKindLParen) {
		// TODO: parse members and member default values
	}

	return
}

// parseBasicDefaultValue parses a default value of a non-object type.
func (parser *ParsingOperation) parseBasicDefaultValue () (
	value Argument,
	err error,
) {
	err = parser.expect(lexer.TokenKindLessThan)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	var arguments []Argument

	defer func () {
		// if we have multiple values, we need to return the full array
		// instead.
		if len(arguments) > 1 {
			// FIXME: i think this is a sign that the tree needs
			// cleaning up.
			value.kind = ArgumentKindArrayDefaultValues
			location := value.location
			value.value = ArrayDefaultValues {
				values: arguments,
			}
			value.location = location
		}
	} ()

	for {
		if parser.token.Is(lexer.TokenKindIndent)   { continue }
		if parser.token.Is(lexer.TokenKindNewline)  { continue }
		if parser.token.Is(lexer.TokenKindGreaterThan) { break }

		value, err = parser.parseArgument()
		if err != nil { return }
		arguments = append(arguments, value)
	}
	return
}

// parseObjectDefaultValue parses default values and new members of an object
// type.
func (parser *ParsingOperation) parseObjectDefaultValue () (
	value   Argument,
	members []TypeMember,
	err error,
) {
	err = parser.expect(lexer.TokenKindLParen)
	if err != nil { return }
	parser.nextToken()
	if err != nil { return }

	for {
		err = parser.expect(lexer.TokenKindDot)
		if err != nil { return }
		parser.nextToken()
		if err != nil { return }
	}
	
	return
}
