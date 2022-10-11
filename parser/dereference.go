package parser

import "git.tebibyte.media/arf/arf/lexer"
// import "git.tebibyte.media/arf/arf/infoerr"

func (parser *ParsingOperation) parseDereference () (
	dereference Dereference,
	err error,
) {
	err = parser.expect(lexer.TokenKindLBrace)
	if err != nil { return }
	dereference.location = parser.token.Location()
	
	// parse the value we are dereferencing
	err = parser.nextToken(validArgumentStartTokens...)
	if err != nil { return }
	dereference.argument, err = parser.parseArgument()
	if err != nil { return }

	// if there is an offset, parse it
	err = parser.expect(lexer.TokenKindUInt, lexer.TokenKindLBrace)
	if err != nil { return }
	if parser.token.Is(lexer.TokenKindUInt) {
		dereference.offset = parser.token.Value().(uint64)
	}
	
	err = parser.nextToken(lexer.TokenKindLBrace)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	return
}
