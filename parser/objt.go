package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

// parseObjtSection parses an object type definition. This allows for structured
// types to be defined, and for member variables to be added and overridden.
func (parser *ParsingOperation) parseObjtSection () (
	section *ObjtSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &ObjtSection {
		location: parser.token.Location(),
		members:  make(map[string] ObjtMember),
	}

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
	section.inherits, err = parser.parseType()
	if err != nil { return }
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	// parse members
	err = parser.parseObjtMembers(section)
	
	if len(section.members) == 0 {
		infoerr.NewError (
			section.location,
			"defining an object with no members",
			infoerr.ErrorKindWarn).Print()
	}
	return
}

func (parser *ParsingOperation) parseObjtMembers (
	into *ObjtSection,
) (
	err     error,
) {
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }

		// get permission
		err = parser.nextToken(lexer.TokenKindPermission)
		if err != nil { return }
		member := ObjtMember { }
		member.permission = parser.token.Value().(types.Permission)

		// get name
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		member.name = parser.token.Value().(string)

		// get type
		err = parser.nextToken(lexer.TokenKindColon)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
		member.what, err = parser.parseType()
		if err != nil { return }
		
		// parse default value
		if parser.token.Is(lexer.TokenKindNewline) {
			err = parser.nextToken()
			if err != nil { return }

			member.defaultValue,
			err = parser.parseInitializationValues(1)
			if err != nil { return }
		} else {
			member.defaultValue, err = parser.parseArgument()
			if err != nil { return }

			err = parser.expect(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
		}

		// add member to object section
		into.members[member.name] = member
	}
}
