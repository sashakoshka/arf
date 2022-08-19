package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
// import "git.tebibyte.media/sashakoshka/arf/infoerr"

func (parser *ParsingOperation) parseTypeSection () (
	section *TypeSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &TypeSection { location: parser.token.Location() }

	// get permission
	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	// get inherited type
	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	section.inherits, err = parser.parseType()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindNewline) {
		err = parser.nextToken()
		if err != nil { return }

		// section.value, err = parser.parseInitializationValues(0)
		// if err != nil { return }
	} else {
		section.defaultValue, err = parser.parseArgument()
		if err != nil { return }

		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
	return
}
