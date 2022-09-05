package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseFunc parses a function section.
func (parser *ParsingOperation) parseFuncSection () (
	section FuncSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section.location = parser.token.Location()

	// get permission
	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	// get arguments
	err = parser.nextToken(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	err = parser.parseFuncArguments(&section)
	if err != nil { return }

	// skip the rest of the section if we are only skimming it
	if parser.skimming {
		section.external = true
		err = parser.skipIndentLevel(1)
		return
	}

	// check to see if the function is external
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != 1         { return }
	err = parser.nextToken()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindName) &&
		parser.token.Value().(string) == "external" {
		
		section.external = true
		err = parser.nextToken(lexer.TokenKindNewline)
		if err != nil { return }
		if err != nil { return }
		err = parser.nextToken()
		
		return
	}

	// if it isn't, backtrack to the start of the line
	parser.previousToken()

	// parse root block
	section.root, err = parser.parseBlock(1)
	if err != nil { return }

	if len(section.root) == 0 {
		infoerr.NewError (section.location,
			"this function has nothing in it",
			infoerr.ErrorKindWarn).Print()
	}
	
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
			err = parser.nextToken(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
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
			
			// get type
			err = parser.nextToken(lexer.TokenKindColon)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			reciever.what, err = parser.parseType()
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
			
			// get type
			err = parser.nextToken(lexer.TokenKindColon)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			input.what, err = parser.parseType()
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
			
			// get type
			err = parser.nextToken(lexer.TokenKindColon)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
			output.what, err = parser.parseType()
			if err != nil { return }
			
			// parse default value
			if parser.token.Is(lexer.TokenKindNewline) {
				err = parser.nextToken()
				if err != nil { return }

				output.value, err =
					parser.parseInitializationValues(1)
				into.outputs = append(into.outputs, output)
				if err != nil { return }
			} else {
				output.value, err =
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
