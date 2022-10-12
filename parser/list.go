package parser

import "git.tebibyte.media/arf/arf/lexer"

// parseList parses a parenthetically delimited list of arguments.
func (parser *parsingOperation) parseList () (list List, err error) {
	list.location = parser.token.Location()

	err = parser.expect(lexer.TokenKindLParen)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	for {
		err = parser.skipWhitespace()
		if err != nil { return }

		// if we have reached the end of the list, stop
		if parser.token.Is(lexer.TokenKindRParen) { break }

		// otherwise, parse argument
		var argument Argument
		argument, err = parser.parseArgument()
		list.arguments = append(list.arguments, argument)
		if err != nil { return }
	}
	
	err = parser.nextToken()
	if err != nil { return }
	
	return
}
