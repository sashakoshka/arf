package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"

// parseData parses a data section.
func (parser *ParsingOperation) parseDataSection () (
	section *DataSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &DataSection { }
	section.setLocation(parser.token.Location())

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
