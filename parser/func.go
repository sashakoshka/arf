package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// validBlockLevelPhraseTokens lists all tokens that are expected when parsing
// a block level phrase.
var validBlockLevelPhraseTokens = append (
	validArgumentStartTokens,
	lexer.TokenKindNewline,
	lexer.TokenKindReturnDirection)
	
// validDelimitedBlockLevelPhraseTokens is like validBlockLevelPhraseTokens, but
// it also includes a right brace token.
var validDelimitedBlockLevelPhraseTokens = append (
	validArgumentStartTokens,
	lexer.TokenKindNewline,
	lexer.TokenKindIndent,
	lexer.TokenKindRBracket,
	lexer.TokenKindReturnDirection)

// parseFunc parses a function section.
func (parser *ParsingOperation) parseFuncSection () (
	section *FuncSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &FuncSection { location: parser.token.Location() }

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
	err = parser.parseFuncArguments(section)
	if err != nil { return }

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

// parseBlock parses an indented block of 
func (parser *ParsingOperation) parseBlock (
	indent int,
) (
	block Block,
	err error,
) {
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != indent    { return }

		var phrase Phrase
		phrase, err = parser.parseBlockLevelPhrase(indent)
		block = append(block, phrase)
		if err != nil { return }
	}
	return
}

// parseBlockLevelPhrase parses a phrase that is not being used as an argument
// to something else. This method is allowed to do things like parse return
// directions, and indented blocks beneath the phrase.
func (parser *ParsingOperation) parseBlockLevelPhrase (
	indent int,
) (
	phrase Phrase,
	err error,
) {
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != indent    { return }
	err = parser.nextToken(validArgumentStartTokens...)
	if err != nil { return }

	expectRightBracket := false
	if parser.token.Is(lexer.TokenKindLBracket) {
		expectRightBracket = true
		err = parser.nextToken()
		if err != nil { return }
	}

	// get command
	err = parser.expect(validArgumentStartTokens...)
	if err != nil { return }
	phrase.command, err = parser.parseArgument()
	if err != nil { return }

	for {
		if expectRightBracket {
			// delimited
			// [someFunc arg1 arg2 arg3] -> someVariable
			err = parser.expect(validDelimitedBlockLevelPhraseTokens...)
			if err != nil { return }

			if parser.token.Is(lexer.TokenKindRBracket) {
				// this is an ending delimiter
				err = parser.nextToken()
				if err != nil { return }
				break
				
			} else if parser.token.Is(lexer.TokenKindNewline) {
				// we are delimited, so we can safely skip
				// newlines
				err = parser.nextToken()
				if err != nil { return }
				continue
				
			} else if parser.token.Is(lexer.TokenKindIndent) {
				// we are delimited, so we can safely skip
				// indents
				err = parser.nextToken()
				if err != nil { return }
				continue
			}
		} else {
			// not delimited
			// someFunc arg1 arg2 arg3 -> someVariable
			err = parser.expect(validBlockLevelPhraseTokens...)
			if err != nil { return }

			if parser.token.Is(lexer.TokenKindReturnDirection) {
				// we've reached a return direction, so that
				// means this is the end of the phrase
				break
			} else if parser.token.Is(lexer.TokenKindNewline) {
				// we've reached the end of the line, so that
				// means this is the end of the phrase.
				break
			}
		}

			
		// this is an argument
		var argument Argument
		argument, err = parser.parseArgument()
		phrase.arguments = append(phrase.arguments, argument)
	}

	// TODO: expect return direction, or newline. then go onto the next
	// line, parsing returnsTo if nescessary.
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	return
}

// parseArgumentLevelPhrase parses a phrase that is being used as an argument to
// something. It is forbidden from using return direction, and it must be
// delimited by brackets.
func (parser *ParsingOperation) parseArgumentLevelPhrase () (
	phrase Phrase,
	err error,
) {
	err = parser.expect(lexer.TokenKindLBracket)
	if err != nil { return }

	// get command
	err = parser.nextToken(validArgumentStartTokens...)
	if err != nil { return }
	phrase.command, err = parser.parseArgument()
	if err != nil { return }

	for {
		// delimited
		// [someFunc arg1 arg2 arg3] -> someVariable
		err = parser.expect(validDelimitedBlockLevelPhraseTokens...)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindRBracket) {
			// this is an ending delimiter
			err = parser.nextToken()
			if err != nil { return }
			return
			
		} else if parser.token.Is(lexer.TokenKindNewline) {
			// we are delimited, so we can safely skip
			// newlines
			err = parser.nextToken()
			if err != nil { return }
			continue
			
		} else if parser.token.Is(lexer.TokenKindIndent) {
			// we are delimited, so we can safely skip
			// indents
			err = parser.nextToken()
			if err != nil { return }
			continue
		}
			
		// this is an argument
		var argument Argument
		argument, err = parser.parseArgument()
		phrase.arguments = append(phrase.arguments, argument)
	}
}
