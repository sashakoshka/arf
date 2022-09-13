package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"

// parseData parses a data section.
func (parser *ParsingOperation) parseDataSection () (
	section DataSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section.location = parser.token.Location()

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

	// skip the rest of the section if we are only skimming it
	if parser.skimming {
		section.external = true
		err = parser.skipIndentLevel(1)
		return
	}

	if parser.token.Is(lexer.TokenKindNewline) {
		err = parser.nextToken()
		if err != nil { return }

		// check if external
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1     { return }
		
		err = parser.nextToken()
		if err != nil { return }
		if parser.token.Is(lexer.TokenKindName) &&
			parser.token.Value().(string) == "external" {
		
			section.external = true
			err = parser.nextToken(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			return
		}

		// otherwise, parse initialization values
		parser.previousToken()
		section.value, err = parser.parseDefaultValues(0)
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
