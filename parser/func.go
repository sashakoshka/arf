package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseFunc parses a function section.
func (parser *ParsingOperation) parseFuncSection () (
	section *FuncSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &FuncSection { location: parser.token.Location() }

	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	err = parser.nextToken(lexer.TokenKindNewline)
	if err != nil { return }
	
	return
}

// parseFuncArguments parses a function's inputs, outputs, and reciever if that
// exists.
func (parser *ParsingOperation) parseFuncArguments (
	into *FuncSection,
) (
	err error,
) {
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) ||
			parser.token.Value().(int) != 1 {
			
			if into.receiver != nil {
				err = parser.token.NewError (
					"func section terminated without a " +
					"separator token",
					infoerr.ErrorKindError)
			}
			return
		}

		// determine whether this is an input, output, or the method
		// reciever
		err = parser.nextToken (
			lexer.TokenKindAt,
			lexer.TokenKindLessThan,
			lexer.TokenKindGreaterThan,
			lexer.TokenKindSeparator)
		if err != nil { return }

		startToken := parser.token
		if parser.token.Is(lexer.TokenKindSeparator) {
			// if we have encountered a separator, that means our
			// work is done here.
			err = parser.expect(lexer.TokenKindNewline)
			return
		}
		
		switch startToken.Kind() {
		case lexer.TokenKindAt:
			reciever := Declaration { }
			reciever.location = parser.token.Location()
			
			// get name
			err = parser.nextToken(lexer.TokenKindName)
			if err != nil { return }
			reciever.name = parser.token.Value().(string)
			err = parser.nextToken()
			if err != nil { return }
			
			if into.receiver != nil {
				err = startToken.NewError (
					"cannot have more than one method " +
					"receiver",
					infoerr.ErrorKindError)
				return
			} else {
				into.receiver = &reciever
			}
			
			err = parser.expect(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			
		case lexer.TokenKindGreaterThan:
			input := Declaration { }
			input.location = parser.token.Location()
			
			// get name
			err = parser.nextToken(lexer.TokenKindName)
			if err != nil { return }
			input.name = parser.token.Value().(string)
			err = parser.nextToken()
			if err != nil { return }

			into.inputs = append(into.inputs, input)
			
			err = parser.expect(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			
		case lexer.TokenKindLessThan:
			output := FuncOutput { }
			output.location = parser.token.Location()
			
			// get name
			err = parser.nextToken(lexer.TokenKindName)
			if err != nil { return }
			output.name = parser.token.Value().(string)
			err = parser.nextToken()
			if err != nil { return }
			
			// parse default value
			if parser.token.Is(lexer.TokenKindNewline) {
				err = parser.nextToken()
				if err != nil { return }

				output.defaultValue, err =
					parser.parseInitializationValues(1)
				into.outputs = append(into.outputs, output)
				if err != nil { return }
			} else {
				output.defaultValue, err =
					parser.parseArgument()
				into.outputs = append(into.outputs, output)
				if err != nil { return }

				err = parser.expect(lexer.TokenKindNewline)
				if err != nil { return }
				err = parser.nextToken()
				if err != nil { return }
			}
		}
	
	}
}
