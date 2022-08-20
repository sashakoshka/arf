package parser

import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

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
	section.root, err = parser.parseTypeNode(0)

	return
}

// parseTypeNode parses a single type definition node recursively.
func (parser *ParsingOperation) parseTypeNode (
	baseIndent int,
) (
	node TypeNode,
	err  error,
) {
	node.children = make(map[string] TypeNode)

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
		
	} else if parser.token.Is(lexer.TokenKindPermission) {

		// child members
		parser.previousToken()
		err = parser.parseTypeNodeChildren(parent)
		
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

// parseTypeNodeChildren parses child type nodes into a parent type node.
func (parser *ParsingOperation) parseTypeNodeChildren (
	parent *TypeNode,
) (
	err error,
) {
	baseIndent := 0
	begin      := true
	
	for {
		// if there is no indent we can just stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { break }
		indent := parser.token.Value().(int)
		
		if begin == true {
			baseIndent = indent 
			begin      = false
		}

		// do not parse any further if the indent has changed
		if indent != baseIndent { break }

		// move on to the beginning of the line, which must contain
		// a type node
		err = parser.nextToken()
		if err != nil { return }
		var child TypeNode
		child, err = parser.parseTypeNode(baseIndent)

		// if the member has already been listed, throw an error
		_, exists := parent.children[child.name]
		if exists {
			err = parser.token.NewError (
				"duplicate member \"" + child.name +
				"\" in object member initialization",
				infoerr.ErrorKindError)
			return
		}

		// store in parent
		parent.children[child.name] = child
	}
	
	return
}
