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

		// section.value, err = parser.parseInitializationValues(0)
		// if err != nil { return }
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
