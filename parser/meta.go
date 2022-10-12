package parser

import "path/filepath"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseMeta parsese the metadata header at the top of an arf file.
func (parser *parsingOperation) parseMeta () (err error) {
	for {
		err = parser.expect (
			lexer.TokenKindName,
			lexer.TokenKindSeparator)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindSeparator) {
			err = parser.nextToken(lexer.TokenKindNewline)
			if err != nil { return }
			
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
			// if import path is relative, get absolute path.
			if value[0] == '.' {
				value = filepath.Join(parser.modulePath, value)
			} else if value[0] != '/' {
				// TODO: get arf import path from an env
				// variable, and default to this if not found.
				// then, search all paths.
				value = filepath.Join (
					"/usr/local/include/arf/",
					value)
			}

			basename  := filepath.Base(value)
			_, exists := parser.tree.requires[basename]

			if exists {
				err = parser.token.NewError (
					"cannot require \"" + basename +
					"\" multiple times",
					infoerr.ErrorKindError)
				return
			}

			parser.tree.requires[basename] = value
		default:
			parser.token.NewError (
				"unrecognized metadata field: " + field,
				infoerr.ErrorKindError)
		}

		err = parser.nextToken(lexer.TokenKindNewline)
		if err != nil { return }
		
		err = parser.nextToken()
		if err != nil { return }
	}
}
