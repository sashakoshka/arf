package parser
// 
// import "git.tebibyte.media/arf/arf/types"
// import "git.tebibyte.media/arf/arf/lexer"
// import "git.tebibyte.media/arf/arf/infoerr"
// 
// // parseObjtSection parses an object type definition. This allows for structured
// // types to be defined, and for member variables to be added and overridden.
// func (parser *ParsingOperation) parseObjtSection () (
	// section ObjtSection,
	// err     error,
// ) {
	// err = parser.expect(lexer.TokenKindName)
	// if err != nil { return }
	// 
	// section.location = parser.token.Location()
// 
	// // get permission
	// err = parser.nextToken(lexer.TokenKindPermission)
	// if err != nil { return }
	// section.permission = parser.token.Value().(types.Permission)
// 
	// // get name
	// err = parser.nextToken(lexer.TokenKindName)
	// if err != nil { return }
	// section.name = parser.token.Value().(string)
// 
	// // parse inherited type
	// err = parser.nextToken(lexer.TokenKindColon)
	// if err != nil { return }
	// err = parser.nextToken()
	// if err != nil { return }
	// section.inherits, err = parser.parseIdentifier()
	// if err != nil { return }
	// err = parser.expect(lexer.TokenKindNewline)
	// if err != nil { return }
	// err = parser.nextToken()
	// if err != nil { return }
// 
	// // parse members
	// err = parser.parseObjtMembers(&section)
	// if err != nil { return }
	// 
	// if len(section.members) == 0 {
		// infoerr.NewError (
			// section.location,
			// "defining an object with no members",
			// infoerr.ErrorKindWarn).Print()
	// }
	// return
// }
// 
// // parseObjtMembers parses a list of members for an object section. Indentation
// // level is assumed.
// func (parser *ParsingOperation) parseObjtMembers (
	// into *ObjtSection,
// ) (
	// err     error,
// ) {
	// for {
		// // if we've left the block, stop parsing
		// if !parser.token.Is(lexer.TokenKindIndent) { return }
		// if parser.token.Value().(int) != 1         { return }
		// 
		// // add member to object section
		// var member ObjtMember
		// member, err = parser.parseObjtMember()
		// into.members = append(into.members, member)
		// if err != nil { return }
	// }
// }
// 
// // parseObjtMember parses a single member of an object section. Indentation
// // level is assumed.
// func (parser *ParsingOperation) parseObjtMember () (
	// member ObjtMember,
	// err    error,
// ) {
	// // get permission
	// err = parser.nextToken(lexer.TokenKindPermission)
	// if err != nil { return }
	// member.permission = parser.token.Value().(types.Permission)
// 
	// // get name
	// err = parser.nextToken(lexer.TokenKindName)
	// if err != nil { return }
	// member.name = parser.token.Value().(string)
// 
	// // get type
	// err = parser.nextToken(lexer.TokenKindColon)
	// if err != nil { return }
	// err = parser.nextToken()
	// if err != nil { return }
	// member.what, err = parser.parseType()
	// if err != nil { return }
// 
	// // if there is a bit width, get it
	// if parser.token.Is(lexer.TokenKindBinaryAnd) {
		// err = parser.nextToken(lexer.TokenKindUInt)
		// if err != nil { return }
		// member.bitWidth = parser.token.Value().(uint64)
		// err = parser.nextToken()
		// if err != nil { return }
	// }
	// 
	// // parse default value
	// if parser.token.Is(lexer.TokenKindNewline) {
		// err = parser.nextToken()
		// if err != nil { return }
// 
		// member.value,
		// err = parser.parseInitializationValues(1)
		// if err != nil { return }
	// } else {
		// member.value, err = parser.parseArgument()
		// if err != nil { return }
// 
		// err = parser.expect(lexer.TokenKindNewline)
		// if err != nil { return }
		// err = parser.nextToken()
		// if err != nil { return }
	// }
// 
	// return
// }
