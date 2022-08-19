package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
// import "git.tebibyte.media/sashakoshka/arf/infoerr"

// parseTypeSection parses a type definition.
func (parser *ParsingOperation) parseTypeSection () (
	section *TypeSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &TypeSection { location: parser.token.Location() }

	// parse root node
	err = parser.nextToken()
	if err != nil { return }
	section.root, err = parser.parseTypeNode(true)

	return
}

// parseTypeNode parses a single type definition node recursively. If isRoot is
// true, the parser will assume the current node is the root node of a type
// section and will not search for an indent.
func (parser *ParsingOperation) parseTypeNode (
	isRoot bool,
) (
	node TypeNode,
	err  error,
) {
	node.children = make(map[string] TypeNode)

	// determine the indent level of this type node. if this is the root
	// node, assume the indent level is zero.
	baseIndent := 0
	if !isRoot {
		err = parser.expect(lexer.TokenKindIndent)
		if err != nil { return }
		baseIndent = parser.token.Value().(int)
		err = parser.nextToken()
		if err != nil { return }
	}

	// get permission
	err = parser.expect(lexer.TokenKindPermission)
	if err != nil { return }
	node.permission = parser.token.Value().(types.Permission)

	// get name
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	node.name = parser.token.Value().(string)

	// get inherited type
	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	node.what, err = parser.parseType()
	if err != nil { return }

	// get value, or child nodes
	if parser.token.Is(lexer.TokenKindNewline) {
		err = parser.nextToken()
		if err != nil { return }

		err = parser.parseTypeNodeBlock(baseIndent, &node)
		if err != nil { return }
	} else {
		node.defaultValue, err = parser.parseArgument()
		if err != nil { return }

		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
	return
}

// parseTypeNodeBlock starts on the line after a type node, and parses what
// could be either an array initialization, an object initialization, or more
// child nodes. It is similar to parseInitializationValues. If none of these
// things were found the parser stays at the beginning of the line and the
// method returns.
func (parser *ParsingOperation) parseTypeNodeBlock (
	baseIndent int,
	parent     *TypeNode,
) (
	err error,
) {
	// check if line is indented one more than baseIndent
	if !parser.token.Is(lexer.TokenKindIndent) { return }
	if parser.token.Value().(int) != baseIndent + 1 { return }

	thingLocation := parser.token.Location()
	
	err = parser.nextToken()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindDot) {

		// object initialization
		parser.previousToken()
		initializationArgument := Argument { location: thingLocation }
		var initializationValues ObjectInitializationValues
		initializationValues, err    = parser.parseObjectInitializationValues()
		initializationArgument.kind  = ArgumentKindObjectInitializationValues
		initializationArgument.value = &initializationValues
		parent.defaultValue = initializationArgument
		
	} else {
	
		// array initialization
		parser.previousToken()
		initializationArgument := Argument { location: thingLocation }
		var initializationValues ArrayInitializationValues
		initializationValues, err    = parser.parseArrayInitializationValues()
		initializationArgument.kind  = ArgumentKindArrayInitializationValues
		initializationArgument.value = &initializationValues
		parent.defaultValue = initializationArgument
	}
	
	return
}
