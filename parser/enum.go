package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

func (parser *ParsingOperation) parseEnumSection () (
	section *EnumSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &EnumSection {
		location: parser.token.Location(),
		members:  make(map[string] Argument),
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
	section.what, err = parser.parseType()
	if err != nil { return }
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	// parse members
	err = parser.parseEnumMembers(section)
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
		
		// get name
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		name := parser.token.Value().(string)
		err = parser.nextToken()
		if err != nil { return }
	
		// parse default value
		var argument Argument
		if parser.token.Is(lexer.TokenKindNewline) {
			err = parser.nextToken()
			if err != nil { return }

			argument, err = parser.parseInitializationValues(1)
			into.members[name] = argument
			if err != nil { return }
		} else {
			argument, err = parser.parseArgument()
			into.members[name] = argument
			if err != nil { return }

			err = parser.expect(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
		}
	}
}
