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
		
		err = parser.nextToken(
			lexer.TokenKindName,
			lexer.TokenKindUInt,
			lexer.TokenKindLParen,
			lexer.TokenKindLessThan)
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
		} else if parser.token.Is(lexer.TokenKindUInt) {
			// parse fixed array length
			what.length = parser.token.Value().(uint64)
			
		} else if parser.token.Is(lexer.TokenKindLessThan) {
			// parse default value
			what.defaultValue, err = parser.parseBasicDefaultValue()
			if err != nil { return }
			
		} else if parser.token.Is(lexer.TokenKindLParen) {
			// parse members and member default values
			what.defaultValue,
			what.members,
			err = parser.parseObjectDefaultValueAndMembers()
			if err != nil { return }
		}
		
		err = parser.nextToken()
		if err != nil { return }
	}

	return
}

// parseBasicDefaultValue parses a default value of a non-object type.
func (parser *ParsingOperation) parseBasicDefaultValue () (
	value Argument,
	err error,
) {
	value.location = parser.token.Location()
	
	err = parser.expect(lexer.TokenKindLessThan)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	var attributes []Argument

	defer func () {
		// if we have multiple values, we need to return the full array
		// instead.
		if len(attributes) > 1 {
			value.kind = ArgumentKindArrayDefaultValues
			value.value = ArrayDefaultValues(attributes)
		}
	} ()

	for {
		if parser.token.Is(lexer.TokenKindIndent)   { continue }
		if parser.token.Is(lexer.TokenKindNewline)  { continue }
		if parser.token.Is(lexer.TokenKindGreaterThan) { break }

		value, err = parser.parseArgument()
		if err != nil { return }
		attributes = append(attributes, value)
	}
	return
}

// parseObjectDefaultValueAndMembers parses default values and new members of an
// object type.
func (parser *ParsingOperation) parseObjectDefaultValueAndMembers () (
	value   Argument,
	members []TypeMember,
	err error,
) {
	value.location = parser.token.Location()

	err = parser.expect(lexer.TokenKindLParen)
	if err != nil { return }
	parser.nextToken()
	if err != nil { return }

	var attributes ObjectDefaultValues

	for {
		if parser.token.Is(lexer.TokenKindIndent)  { continue }
		if parser.token.Is(lexer.TokenKindNewline) { continue }
		if parser.token.Is(lexer.TokenKindRParen)  { break }
		
		err = parser.expect(lexer.TokenKindDot)
		if err != nil { return }
		parser.nextToken(lexer.TokenKindName, lexer.TokenKindPermission)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindName) {
			// parsing a defalut value for an inherited member
			var memberName  string
			var memberValue Argument

			memberName,
			memberValue, err = parser.parseObjectMemberDefinition()
			if err != nil { return }
			
			if value.kind == ArgumentKindNil {
				// create default value map if it doesn't
				// already exist
				value.kind = ArgumentKindObjectDefaultValues
				attributes = make(ObjectDefaultValues)
				value.value = attributes
			}

			// TODO: error on duplicate
			if memberValue.kind != ArgumentKindNil {
				attributes[memberName] = memberValue
			}
			
		} else if parser.token.Is(lexer.TokenKindPermission) {
			// parsing a member declaration
			var member TypeMember
			member,
			err = parser.parseObjectMemberDeclaration()

			members = append(members, member)
		}
	}
	
	return
}

// parseObjectDefaultValue parses member default values only, and will throw an
// error when it encounteres a new member definition.
func (parser *ParsingOperation) parseObjectDefaultValue () (
	value Argument,
	err error,
) {
	value.location = parser.token.Location()

	err = parser.expect(lexer.TokenKindLParen)
	if err != nil { return }
	parser.nextToken()
	if err != nil { return }

	var attributes ObjectDefaultValues

	for {
		if parser.token.Is(lexer.TokenKindIndent)  { continue }
		if parser.token.Is(lexer.TokenKindNewline) { continue }
		if parser.token.Is(lexer.TokenKindRParen)  { break }
		
		err = parser.expect(lexer.TokenKindDot)
		if err != nil { return }
		parser.nextToken(lexer.TokenKindName)
		if err != nil { return }

		if value.kind == ArgumentKindNil {
			value.kind = ArgumentKindObjectDefaultValues
			attributes = make(ObjectDefaultValues)
			value.value = attributes
		}

		var memberName  string
		var memberValue Argument
		memberName,
		memberValue, err = parser.parseObjectMemberDefinition()

		attributes[memberName] = memberValue
	}
	
	return
}

// .ro name:Type:qualifier:<value>

// parseObjectMemberDeclaration parses a new default value for an inherited
// member.
func (parser *ParsingOperation) parseObjectMemberDefinition () (
	name  string,
	value Argument,
	err   error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	return
}

// .name:<value>

// parseObjectMemberDeclaration parses an object member declaration, and its
// default value if it exists.
func (parser *ParsingOperation) parseObjectMemberDeclaration () (
	member TypeMember,
	err error,
) {
	err = parser.expect(lexer.TokenKindPermission)
	if err != nil { return }
	
	return
}
