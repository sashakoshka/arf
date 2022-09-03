package parser

import "git.tebibyte.media/arf/arf/lexer"

// parseBlock parses an indented block of phrases
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
	err = parser.nextToken(validPhraseStartTokens...)
	if err != nil { return }

	expectRightBracket := false
	if parser.token.Is(lexer.TokenKindLBracket) {
		expectRightBracket = true
		err = parser.nextToken()
		if err != nil { return }
	}

	// get command
	err = parser.expect(validPhraseStartTokens...)
	if err != nil { return }
	if isTokenOperator(parser.token) {
		phrase.command, err = parser.parseOperatorArgument()
		if err != nil { return }
	} else {
		phrase.command, err = parser.parseArgument()
		if err != nil { return }
	}

	for {
		if expectRightBracket {
			// delimited
			// [someFunc arg1 arg2 arg3] -> someVariable
			err = parser.expect(validDelimitedPhraseTokens...)
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
	err = parser.nextToken(validPhraseStartTokens...)
	if err != nil { return }
	if isTokenOperator(parser.token) {
		phrase.command, err = parser.parseOperatorArgument()
		if err != nil { return }
	} else {
		phrase.command, err = parser.parseArgument()
		if err != nil { return }
	}

	for {
		// delimited
		// [someFunc arg1 arg2 arg3] -> someVariable
		err = parser.expect(validDelimitedPhraseTokens...)
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

// parseOperatorArgument parses an operator argument. This is only used to parse
// operator phrase commands.
func (parser *ParsingOperation) parseOperatorArgument () (
	operator Argument,
	err error,
) {
	err = parser.expect(operatorTokens...)
	if err != nil { return }

	operator.location = parser.token.Location()
	operator.kind     = ArgumentKindOperator
	operator.value    = parser.token.Kind()
	
	err = parser.nextToken()
	return
}
