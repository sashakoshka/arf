package parser

import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/lexer"

// parseData parses a data section.
func (parser *ParsingOperation) parseDataSection () (
	section *DataSection,
	err     error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	
	section = &DataSection { location: parser.token.Location() }

	err = parser.nextToken(lexer.TokenKindPermission)
	if err != nil { return }
	section.permission = parser.token.Value().(types.Permission)

	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }
	section.name = parser.token.Value().(string)

	err = parser.nextToken(lexer.TokenKindColon)
	if err != nil { return }
	err = parser.nextToken()
	if err != nil { return }
	section.what, err = parser.parseType()
	if err != nil { return }

	if parser.token.Is(lexer.TokenKindNewline) {
		err = parser.nextToken()
		if err != nil { return }

		section.value, err = parser.parseInitializationValues(0)
		if err != nil { return }
	} else {
		section.value, err = parser.parseArgument()
		if err != nil { return }

		err = parser.expect(lexer.TokenKindNewline)
		if err != nil { return }
		err = parser.nextToken()
		if err != nil { return }
	}
	return
}

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

	if parser.token.Is(lexer.TokenKindDot) {
	
		var initializationValues ObjectInitializationValues
		initializationValues, err    = parser.parseObjectInitializationValues()
		initializationArgument.kind  = ArgumentKindObjectInitializationValues
		initializationArgument.value = &initializationValues
		
	} else {
	
		var initializationValues ArrayInitializationValues
		initializationValues, err    = parser.parseArrayInitializationValues()
		initializationArgument.kind  = ArgumentKindArrayInitializationValues
		initializationArgument.value = &initializationValues
		
	}
	
	return
}

// parseObjectInitializationValues parses a list of object initialization
// values until the indentation level drops.
func (parser *ParsingOperation) parseObjectInitializationValues () (
	values ObjectInitializationValues,
	err    error,
) {
	// println("PARSING")
	// defer println("DONE\n")
	// values.attributes = make(map[string] Argument)
	// values.location   = parser.token.Location()
// 
	// for {
		// // get attribute name and value
		// err = parser.nextToken(lexer.TokenKindName)
		// if err != nil { return }
		// name := parser.token.Value().(string)
		// println("  name:", name)
		// 
		// _, exists := values.attributes[name]
		// if exists {
			// err = parser.token.NewError (
				// "duplicate member \"" + name + "\" in object " +
				// "member initialization",
				// file.ErrorKindError)
			// return
		// }
		// 
		// err = parser.nextToken()
		// if err != nil { return }
// 
		// println("  parsing value argument")
		// // parse value argument
		// var value Argument
		// if parser.token.Is(lexer.TokenKindNewline) {
			// println("    is newline")
			// // if there is none on this line, expect complex
			// // initialization below
// 
			// possibleErrorLocation := parser.token.Location()
			// err = parser.nextToken(lexer.TokenKindIndent)
			// if err != nil { return }
// 
			// value, err = parser.parseInitializationValues(indent)
			// if err != nil { return }
// 
			// // TODO: this doesn't seem to produce an error at the
			// // correct location.
			// if value.value == nil {
				// err = possibleErrorLocation.NewError (
					// "empty initialization value",
					// file.ErrorKindError)
				// return
			// }
			// 
		// } else {
			// println("    is not newline")
			// value, err = parser.parseArgument()
			// if err != nil { return }
			// err = parser.expect(lexer.TokenKindNewline)
			// if err != nil { return }
		// }
// 
		// // store in object
		// values.attributes[name] = value
// 
		// // if indent drops, or does something strange, stop parsing
		// err = parser.nextToken()
		// if err != nil { return }
		// if !parser.token.Is(lexer.TokenKindIndent) { break }
		// if parser.token.Value().(int) != indent { break }
// 
		// // the next line must start with a dot
		// err = parser.nextToken(lexer.TokenKindDot)
		// if err != nil { return }
	// }
	return
}

// parseArrayInitializationValues parses a list of array initialization values
// until the indentation lexel drops.
func (parser *ParsingOperation) parseArrayInitializationValues () (
	values ArrayInitializationValues,
	err    error,
) {
	baseIndent := 0
	begin      := true
	
	for {
		// if there is no indent we can just stop parsing
		if !parser.token.Is(lexer.TokenKindIndent) { break}
		indent := parser.token.Value().(int)
		
		if begin == true {
			values.location = parser.token.Location()
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
			values.values = append(values.values, argument)
		}
	}

	return
}

// parseType parses a type notation of the form Name, {Name}, etc.
func (parser *ParsingOperation) parseType () (what Type, err error) {
	err = parser.expect(lexer.TokenKindName, lexer.TokenKindLBrace)
	if err != nil { return }
	what.location = parser.token.Location()

	if parser.token.Is(lexer.TokenKindLBrace) {
		what.kind = TypeKindPointer

		err = parser.nextToken()
		if err != nil { return }
	
		var points Type
		points, err = parser.parseType()
		if err != nil { return }
		what.points = &points

		err = parser.expect (
			lexer.TokenKindUInt,
			lexer.TokenKindRBrace,
			lexer.TokenKindElipsis)
		if err != nil { return }

		if parser.token.Is(lexer.TokenKindUInt) {
			what.kind = TypeKindArray
		
			what.length = parser.token.Value().(uint64)
		
			err = parser.nextToken(lexer.TokenKindRBrace)
			if err != nil { return }
		} else if parser.token.Is(lexer.TokenKindElipsis) {
			what.kind = TypeKindArray
		
			err = parser.nextToken(lexer.TokenKindRBrace)
			if err != nil { return }
		}

		err = parser.nextToken()
		if err != nil { return }
	} else {
		what.name, err = parser.parseIdentifier()
		if err != nil { return }
	}

	if parser.token.Is(lexer.TokenKindColon) {
		err = parser.nextToken(lexer.TokenKindName)
		if err != nil { return }

		qualifier := parser.token.Value().(string)
		switch qualifier {
		case "mut":
			what.mutable = true
		default:
			err = parser.token.NewError (
				"unknown type qualifier \"" + qualifier + "\"",
				file.ErrorKindError)
			return
		}
		
		err = parser.nextToken()
		if err != nil { return }
	}

	return
}

// parseIdentifier parses an identifier made out of dot separated names.
func (parser *ParsingOperation) parseIdentifier () (
	identifier Identifier,
	err        error,
) {
	err = parser.expect(lexer.TokenKindName)
	if err != nil { return }
	identifier.location = parser.token.Location()

	for {
		// TODO: eat up newlines and tabs after the dot, but not before
		// it.
		if !parser.token.Is(lexer.TokenKindName) { break }

		identifier.trail = append (
			identifier.trail,
			parser.token.Value().(string))

		err = parser.nextToken()
		if err != nil { return }
		
		if !parser.token.Is(lexer.TokenKindDot) { break }
	}

	return
}
