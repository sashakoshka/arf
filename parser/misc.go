package parser

import "git.tebibyte.media/arf/arf/lexer"

// parseIdentifier parses an identifier made out of dot separated names.
func (parser *ParsingOperation) parseIdentifier () (
	identifier Identifier,
	err        error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	identifier.location = parser.token.Location()

	for {
		if !parser.token.Is(lexer.TokenKindName) { break }

		identifier.trail = append (
			identifier.trail,
			parser.token.Value().(string))

		err = parser.nextToken()
		if err != nil { return }
		
		if !parser.token.Is(lexer.TokenKindDot) { break }

		err = parser.nextToken()
		if err != nil { return }

		// allow the identifier to continue on to the next line if there
		// is a line break right after the dot
		for parser.token.Is(lexer.TokenKindNewline) ||
			parser.token.Is(lexer.TokenKindIndent) {

			err = parser.nextToken()
			if err != nil { return }
		}
	}

	return
}
