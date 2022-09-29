package parser

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parseFaceSection parses an interface section.
func (parser *ParsingOperation) parseFaceSection () (
	section FaceSection,
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

	// parse inherited interface
	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.inherits, err = parser.parseIdentifier()
	if err != nil { return }
	err = parser.expect(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != 1         { return }

	err = parser.nextToken (
		lexer.TokenKindName,
		lexer.TokenKindGreaterThan,
		lexer.TokenKindLessThan)
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindName) {
		// parse type interface
		section.kind = FaceKindType
		parser.previousToken()
		section.behaviors, err = parser.parseFaceBehaviors()
		if err != nil { return }
	} else {
		// parse function interface
		section.kind = FaceKindFunc
		parser.previousToken()
		section.inputs,
		section.outputs, err = parser.parseFaceBehaviorArguments(1)
		if err != nil { return }
	}

	return
}

// parseFaceBehaviors parses a list of interface behaviors for an object
// interface.
func (parser *ParsingOperation) parseFaceBehaviors () (
	behaviors map[string] FaceBehavior,
	err error,
) {
	// parse members
	behaviors = make(map[string] FaceBehavior)
	for {
		// if we've left the block, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != 1         { return }
		
		err = parser.nextToken(lexer.TokenKindName)
		behaviorBeginning := parser.token.Location()
		if err != nil { return }

		// parse behavior
		var behavior FaceBehavior
		behavior, err = parser.parseFaceBehavior(1)

		// add to section
		_, exists := behaviors[behavior.name]
		if exists {
			err = infoerr.NewError (
				behaviorBeginning,
				"multiple behaviors named " + behavior.name +
				" in this interface",
				infoerr.ErrorKindError)
			return
		}
		behaviors[behavior.name] = behavior
		
		if err != nil { return }
	}
}

// parseFaceBehavior parses a single interface behavior.
func (parser *ParsingOperation) parseFaceBehavior (
	indent int,
) (
	behavior FaceBehavior,
	err      error,
) {
	// get name
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	behavior.name = parser.token.Value().(string)

	err = parser.nextToken(lexer.TokenKindNewline)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }

	behavior.inputs,
	behavior.outputs,
	err = parser.parseFaceBehaviorArguments(indent + 1)
	if err != nil { return }

	return
}

func (parser *ParsingOperation) parseFaceBehaviorArguments (
	indent int,
) (
	inputs  []Declaration,
	outputs []Declaration,
	err error,
) {
	
	for {
		// if we've left the behavior, stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { return }
		if parser.token.Value().(int) != indent    { return }

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

		if kind == lexer.TokenKindGreaterThan {
			inputs = append(inputs, declaration)
		} else {
			outputs = append(outputs, declaration)
		}
		
		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
}
