package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// TODO:
// (parser *ParsingOperation) parseDefaultValues

// (parser *ParsingOperation) parseDefaultMemberValues (return tree of new members and a tree of member values)
// (parser *ParsingOperation) parseDefaultArrayValues

// (parser *ParsingOperation) parseDefaultMemberValue
// (parser *ParsingOperation) parseMemberDeclaration

// parseInitializationValues starts on the line after a data section, or a set
// phrase. It checks for an indent greater than the indent of the aforementioned
// data section or set phrase (passed through baseIndent), and if there is,
// it parses initialization values.
func (parser *ParsingOperation) parseInitializationValues (
	baseIndent int,
) (
	initializationArgument Argument,
	err error,
) {
	// check if line is indented one more than baseIndent
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != baseIndent + 1 { return }

	initializationArgument.location = parser.token.Location()
	
	err = parser.nextToken()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindDot) {

		// object initialization
		parser.previousToken()
		var initializationValues ObjectInitializationValues
		initializationValues, err    = parser.parseObjectInitializationValues()
		initializationArgument.kind  = ArgumentKindObjectInitializationValues
		initializationArgument.value = initializationValues
		
	} else {
	
		// array initialization
		parser.previousToken()
		var initializationValues ArrayInitializationValues
		initializationValues, err    = parser.parseArrayInitializationValues()
		initializationArgument.kind  = ArgumentKindArrayInitializationValues
		initializationArgument.value = initializationValues
	}
	
	return
}

// parseObjectInitializationValues parses a list of object initialization
// values until the indentation level drops.
func (parser *ParsingOperation) parseObjectInitializationValues () (
	initializationValues ObjectInitializationValues,
	err                  error,
) {
	initializationValues.attributes = make(map[string] Argument)

	baseIndent := 0
	begin      := true
	
	for {
		// if there is no indent we can just stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { break}
		indent := parser.token.Value().(int)
		
		if begin == true {
			initializationValues.location = parser.token.Location()
			baseIndent = indent 
			begin      = false
		}

		// do not parse any further if the indent has changed
		if indent != baseIndent { break }

		// move on to the beginning of the line, which must contain
		// a member initialization value
		err = parser.nextToken(lexer.TokenKindDot)
		if err != nil { return }
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		name := parser.token.Value().(string)

		// if the member has already been listed, throw an error
		_, exists := initializationValues.attributes[name]
		if exists {
			err = parser.token.NewError (
				"duplicate member \"" + name + "\" in object " +
				"member initialization",
				infoerr.ErrorKindError)
			return
		}

		// parse the argument determining the member initialization
		// value
		err = parser.nextToken()
		if err != nil { return }
		var value Argument
		if parser.token.Is(lexer.TokenKindNewline) {
		
			// recurse
			err = parser.nextToken(lexer.TokenKindIndent)
			if err != nil { return }
			
			value, err = parser.parseInitializationValues(baseIndent)
			initializationValues.attributes[name] = value
			if err != nil { return }
			
		} else {

			// parse as normal argument
			value, err = parser.parseArgument()
			initializationValues.attributes[name] = value
			if err != nil { return }
			
			err = parser.expect(lexer.TokenKindNewline)
			if err != nil { return }
			err = parser.nextToken()
			if err != nil { return }
		}
	}
	
	return
}

// parseArrayInitializationValues parses a list of array initialization values
// until the indentation lexel drops.
func (parser *ParsingOperation) parseArrayInitializationValues () (
	initializationValues ArrayInitializationValues,
	err                  error,
) {
	baseIndent := 0
	begin      := true
	
	for {
		// if there is no indent we can just stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { break}
		indent := parser.token.Value().(int)
		
		if begin == true {
			initializationValues.location = parser.token.Location()
			baseIndent = indent 
			begin      = false
		}

		// do not parse any further if the indent has changed
		if indent != baseIndent { break }

		// move on to the beginning of the line, which must contain
		// arguments
		err = parser.nextToken(validArgumentStartTokens...)
		if err != nil { return }

		for {
			// stop parsing this line and go on to the next if a
			// newline token is encountered
			if parser.token.Is(lexer.TokenKindNewline) {
				err = parser.nextToken()
				if err != nil { return }
				break
			}

			// otherwise, parse the argument
			var argument Argument
			argument, err = parser.parseArgument()
			if err != nil { return }
			initializationValues.values = append (
				initializationValues.values,
				argument)
		}
	}

	return
}
