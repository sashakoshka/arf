package parser

import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/lexer"

// parseMeta parsese the metadata header at the top of an arf file.
func (parser *ParsingOperation) parseMeta () (err error) {
	for {
		err = parser.expect (
			lexer.TokenKindName,
			lexer.TokenKindSeparator)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindSeparator) {
			err = parser.nextToken()
			return
		}

		field := parser.token.Value().(string)

		err = parser.nextToken(lexer.TokenKindString)
		if err != nil { return }
		value := parser.token.Value().(string)
		
		switch field {
		case "author":
			parser.tree.author = value
		case "license":
			parser.tree.license = value
		case "require":
			parser.tree.requires = append(parser.tree.requires, value)
		default:
			parser.token.NewError (
				"unrecognized metadata field: " + field,
				file.ErrorKindError)
		}

		err = parser.nextToken(lexer.TokenKindNewline)
		if err != nil { return }
		
		err = parser.nextToken()
		if err != nil { return }
	}
}
