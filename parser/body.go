package parser

import "git.tebibyte.media/sashakoshka/arf/lexer"

// parse body parses the body of an arf file, after the metadata header.
func (parser *ParsingOperation) parseBody () (err error) {
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }

	switch parser.token.Value().(string) {
	case "data":
	case "type":
	case "func":
	case "face":
	}

	return
}
