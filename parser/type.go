package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"

// parseTypeSection parses a type definition. It can inherit from other types,
// and define new members on them.
func (parser *parsingOperation) parseTypeSection () (
	section TypeSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section.location = parser.token.Location()

	// get permission
	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	// parse inherited type
	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	section.what, err = parser.parseType()
	if err != nil { return }

	// see if value exists
	if parser.token.Is(lexer.TokenKindNewline) {
		parser.nextToken()
		// if we have exited the section, return
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }
		
		err = parser.nextToken()
		if err != nil { return }
	}

	// if we have not encountered members, get value and return.
	if !parser.token.Is(lexer.TokenKindPermission) {
		section.argument, err = parser.parseArgument()
		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		
		err = parser.nextToken()
		if err != nil { return }

		return
	}

	parser.previousToken()

	for {
		// if we have exited the section, return
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }
		
		err = parser.nextToken(lexer.TokenKindPermission)
		if err != nil { return }
		var member TypeSectionMember
		member, err = parser.parseTypeSectionMember()
		section.members = append(section.members, member)
		if err != nil { return }
	}
}

// parseTypeSectionMember parses a type section member variable.
func (parser *parsingOperation) parseTypeSectionMember () (
	member TypeSectionMember,
	err error,
) {
	// get permission
	err = parser.expect(lexer.TokenKindPermission)
	if err != nil { return }
	member.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	member.name = parser.token.Value().(string)

	// if there is a type, get it
	err = parser.nextToken()
	if err != nil { return }
	if parser.token.Is(lexer.TokenKindColon) {
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		member.what, err = parser.parseType()
		if err != nil { return }
	}

	// see if value exists
	if parser.token.Is(lexer.TokenKindNewline) {
		parser.nextToken()
		// if we have exited the member, return
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 2         { return }
		
		err = parser.nextToken()
		if err != nil { return }
	}

	// if default value exists, get it
	if !parser.token.Is(lexer.TokenKindBinaryAnd) {
		member.argument, err = parser.parseArgument()
	}

	// if there is a bit width specifier, get it
	if parser.token.Is(lexer.TokenKindBinaryAnd) {
		err = parser.nextToken(lexer.TokenKindUInt)
		if err != nil { return }
		member.bitWidth = parser.token.Value().(uint64)
		
		err = parser.nextToken()
		if err != nil { return }
	}
	
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	
	err = parser.nextToken()
	if err != nil { return }

	return
}
