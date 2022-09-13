package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseEnumSection parses an enumerated type section.
func (parser *ParsingOperation) parseEnumSection () (
	section EnumSection,
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
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	// parse members
	err = parser.parseEnumMembers(&section)
	if err != nil { return }
	
	if len(section.members) == 0 {
		infoerr.NewError (
			section.location,
			"defining an enum with no members",
			infoerr.ErrorKindWarn).Print()
	}
	return
}

// parseEnumMembers parses a list of members for an enum section. Indentation
// level is assumed.
func (parser *ParsingOperation) parseEnumMembers (
	into *EnumSection,
) (
	err error,
) {
	
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }

		member := EnumMember { }
		
		// get name
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		member.location = parser.token.Location()
		member.name = parser.token.Value().(string)
		err = parser.nextToken()
		if err != nil { return }
	
		// parse default value
		member.value, err = parser.parseArgument()
		into.members = append(into.members, member)
		if err != nil { return }

		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
}
