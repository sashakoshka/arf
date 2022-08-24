package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

// parseFaceSection parses an interface section.
func (parser *ParsingOperation) parseFaceSection () (
	section *FaceSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &FaceSection {
		location: parser.token.Location(),
		behaviors:  make(map[string] FaceBehavior),
	}

	// get permission
	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	// parse inherited interface
	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.inherits = parser.token.Value().(string)
	if err != nil { return }
	err = parser.nextToken(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	// parse members
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }

		// parse behavior
		behaviorBeginning := parser.token.Location()
		var behavior FaceBehavior
		behavior, err = parser.parseFaceBehavior()

		// add to section
		_, exists := section.behaviors[behavior.name]
		if exists {
			err = infoerr.NewError (
				behaviorBeginning,
				"multiple behaviors named " + behavior.name +
				" in this interface",
				infoerr.ErrorKindError)
			return
		}
		section.behaviors[behavior.name] = behavior
		
		if err != nil { return }
	}
	return
}

// parseFaceBehavior parses a single interface behavior. Indentation level is
// assumed.
func (parser *ParsingOperation) parseFaceBehavior () (
	behavior FaceBehavior,
	err      error,
) {
	err = parser.expect(lexer.TokenKindIndent)
	if err != nil { return }

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	behavior.name = parser.token.Value().(string)

	err = parser.nextToken(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 2         { return }

		// get preceding symbol
		err = parser.nextToken (
			lexer.TokenKindGreaterThan,
			lexer.TokenKindLessThan)
		if err != nil { return }
		kind := parser.token.Kind()

		var declaration Declaration

		// get name
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }
		declaration.name = parser.token.Value().(string)

		// parse inherited type
		err = parser.nextToken(lexer.TokenKindColon)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
		declaration.what, err = parser.parseType()
		if err != nil { return }
		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }

		if kind == lexer.TokenKindGreaterThan {
			behavior.inputs = append (
				behavior.inputs,
				declaration)
		} else {
			behavior.outputs = append (
				behavior.outputs,
				declaration)
		}
	}
	
	return
}
