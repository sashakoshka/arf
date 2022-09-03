package parser

import "git.tebibyte.media/arf/arf/lexer"

// operatorTokens lists all symbolic tokens that can act as a command to a
// phrase.
var operatorTokens = []lexer.TokenKind {
        lexer.TokenKindColon,
        lexer.TokenKindPlus,
        lexer.TokenKindMinus,
        lexer.TokenKindIncrement,
        lexer.TokenKindDecrement,
        lexer.TokenKindAsterisk,
        lexer.TokenKindSlash,
        lexer.TokenKindExclamation,
        lexer.TokenKindPercent,
        lexer.TokenKindPercentAssignment,
        lexer.TokenKindTilde,
        lexer.TokenKindTildeAssignment,
        lexer.TokenKindEqualTo,
        lexer.TokenKindNotEqualTo,
        lexer.TokenKindLessThanEqualTo,
        lexer.TokenKindLessThan,
        lexer.TokenKindLShift,
        lexer.TokenKindLShiftAssignment,
        lexer.TokenKindGreaterThan,
        lexer.TokenKindGreaterThanEqualTo,
        lexer.TokenKindRShift,
        lexer.TokenKindRShiftAssignment,
        lexer.TokenKindBinaryOr,
        lexer.TokenKindBinaryOrAssignment,
        lexer.TokenKindLogicalOr,
        lexer.TokenKindBinaryAnd,
        lexer.TokenKindBinaryAndAssignment,
        lexer.TokenKindLogicalAnd,
        lexer.TokenKindBinaryXor,
        lexer.TokenKindBinaryXorAssignment,
}

// isTokenOperator returns whether or not the token is an operator token.
func isTokenOperator (token lexer.Token) (isOperator bool) {
	for _, kind := range operatorTokens {
		if token.Is(kind) {
			isOperator = true
			return
		}
	}

	return
}

// validPhraseStartTokens lists all tokens that are expected when parsing the
// first part of a phrase.
var validPhraseStartTokens = append (
	operatorTokens,
	lexer.TokenKindLBracket,
	lexer.TokenKindName,
	lexer.TokenKindString)
	
// validBlockLevelPhraseTokens lists all tokens that are expected when parsing
// a block level phrase.
var validBlockLevelPhraseTokens = append (
	validArgumentStartTokens,
	lexer.TokenKindNewline,
	lexer.TokenKindReturnDirection)
	
// validDelimitedPhraseTokens is like validBlockLevelPhraseTokens, but it also
// includes a right brace token.
var validDelimitedPhraseTokens = append (
	validArgumentStartTokens,
	lexer.TokenKindNewline,
	lexer.TokenKindIndent,
	lexer.TokenKindRBracket,
	lexer.TokenKindReturnDirection)

// controlFlowNames contains a list of all command names that must have a block
// underneath them.
var controlFlowNames = []string {
	"if", "else", "elseif",
	"for", "while",
	"defer",
}

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

	// expect newline or return direction
	err = parser.expect (
		lexer.TokenKindNewline,
		lexer.TokenKindReturnDirection)
	if err != nil { return }
	expectReturnDirection := parser.token.Is(lexer.TokenKindReturnDirection)

	// if we have hit a return direction, parse it...
	if expectReturnDirection {
		err = parser.nextToken()
		if err != nil { return }
		
		for {
			err = parser.expect (
				lexer.TokenKindNewline,
				lexer.TokenKindName)
			if err != nil { return }
			// ...until we hit a newline
			if parser.token.Is(lexer.TokenKindNewline) { break }

			var returnTo Argument
			returnTo, err = parser.parseArgument()
			if err != nil { return }

			phrase.returnsTo = append(phrase.returnsTo, returnTo)
		}
	}
	
	err = parser.nextToken()
	if err != nil { return }

	// if this is a control flow statement, parse block under it
	if phrase.command.kind == ArgumentKindIdentifier {
		// perhaps it is a normal control flow statement?
		command := phrase.command.value.(Identifier)
		if len(command.trail) != 1 { return }
		
		isControlFlow := false
		for _, name := range controlFlowNames {
			if command.trail[0] == name {
				isControlFlow = true
				break
			}
		}

		if !isControlFlow { return }
		
	} else if phrase.command.kind == ArgumentKindOperator {
		// perhaps it is a switch case?
		command := phrase.command.value.(lexer.TokenKind)
		if command != lexer.TokenKindColon { return }
	
	} else {
		return
	}

	// if it is any of those, parse the block under it
	phrase.block, err = parser.parseBlock(indent + 1)

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
