package parser

import "github.com/sashakoshka/arf/lexer"

/* parseBodyFunction parses a function section.
 */
func (parser *Parser) parseBodyFunction (
        skim bool,
) (
        section *Function,
        err error,
) {
        section = &Function {
                where:   parser.embedPosition(),
                root:   &Block {
                        variables: make(map[string] *Variable),
                },
        }

        if !parser.expect(lexer.TokenKindPermission) {
                 return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)

        // if we are skimming and other modules don't have access to this, don't
        // even bother parsing the argument and stuff
        if (skim && section.modeExternal == ModeDeny) {
                section.external = true
                return nil, parser.skipBodySection()
        }

        parser.nextToken()
        if !parser.expect(lexer.TokenKindName) {
                 return nil, parser.skipBodySection()
        }

        section.name = parser.token.StringValue

        parser.nextToken()
        if !parser.expect() { return nil, parser.skipBodySection() }
                
        parser.nextLine()
        if parser.endOfFile() || parser.line.Indent == 0 { return }

        // function arguments
        inHead := true
        for inHead {
                if !parser.expect (
                        lexer.TokenKindSeparator,
                        lexer.TokenKindSymbol,
                ) { return nil, parser.skipBodySection() }

                if parser.token.Kind == lexer.TokenKindSymbol {
                        err = parser.parseBodyFunctionArgumentFor(section)
                        if err != nil { return }
                }
                
                inHead = parser.token.Kind != lexer.TokenKindSeparator

                parser.nextLine()
                if parser.endOfFile() || parser.line.Indent == 0 { return }
        }

        // if we are skimming the file, skip over the function content
        if (skim) {
                section.external = true
                return section, parser.skipBodySection()
        }

        isExternal :=
                parser.token.Kind == lexer.TokenKindName &&
                parser.token.StringValue == "external"
        
        // if the function is external, skip over it.
        if (isExternal) {
                parser.nextToken()
                if parser.token.Kind != lexer.TokenKindNone {
                        parser.printError (
                                parser.token.Column,
                                "nothing should come after external")
                        return nil, nil
                }

                section.external = true
                return section, parser.skipBodySection()
        }

        // function body
        _, err = parser.parseBodyFunctionBlock(0, section.root)
        return
}

/* parseBodyFunctionArgumentFor parses a function argument for the specified
 * function.
 */
func (parser *Parser) parseBodyFunctionArgumentFor (
        section *Function,
) (
        err error,
) {
        switch parser.token.StringValue {
        case "@":
                self := &Variable { where: parser.embedPosition() }
                
                self.name,
                self.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }
                
                if self.what.points == nil {
                        parser.printError (
                                parser.token.Column,
                                "method reciever must be a",
                                "pointer")
                        break
                }
                
                if self.what.mutable {
                        parser.printError (
                                parser.token.Column,
                                "method reciever cannot be",
                                "mutable")
                        break
                }

                if len(self.what.name.trail) > 1 {
                        parser.printError (
                                parser.token.Column,
                                "cannot use member selection in method",
                                "reciever type, type name cannot have dots in",
                                "it")
                }
                
                // add self to function
                if section.root.addVariable(self) {
                        section.self = self.name
                        section.selfType = self.what.points.name.trail[0]
                        section.isMember = true
                } else {
                        parser.printError (
                                parser.token.Column,
                                "a variable with the name", self.name, "is",
                                "already defined in this function")
                }
                break

        case ">":
                input := &Variable { where: parser.embedPosition() }
                
                input.name,
                input.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }

                if input.what.mutable {
                        parser.printError (
                                parser.token.Column,
                                "function arguments cannot be",
                                "mutable")
                        break
                }

                if parser.endOfLine() { break}

                // get default value for input
                input.value,
                _, err = parser.parseDefaultValues(1)
                if err != nil { return err }

                // add input to function
                if section.root.addVariable(input) {
                        section.inputs = append(section.inputs, input.name)
                } else {
                        parser.printError (
                                parser.token.Column,
                                "a variable with the name", input.name, "is",
                                "already defined in this function")
                }
                break
        
        case "<":
                output := &Variable { where: parser.embedPosition() }
                
                output.name,
                output.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }

                if output.what.mutable {
                        parser.printWarning (
                                parser.token.Column,
                                "immutable output, this is useless. consider",
                                "marking as :mut")
                        break
                }

                if parser.endOfLine() { break}

                // get default value for output
                output.value,
                _, err = parser.parseDefaultValues(1)
                if err != nil { return err }

                // add output to function
                if section.root.addVariable(output) {
                        section.outputs = append(section.outputs, output.name)
                } else {
                        parser.printError (
                                parser.token.Column,
                                "a variable with the name", output.name, "is",
                                "already defined in this function")
                }
                break

        default:
                parser.printError (
                        parser.token.Column,
                        "unknown argument type symbol '" +
                        parser.token.StringValue + "',",
                        "use either '@', '>', or '<'")
                break
        }

        return nil
}

/* parseBodyFunctionBlock parses a block of function calls. This is done
 * recursively, so it will also parse sub-blocks. If preExisting is non-nil,
 * this function will parse into it.
 */
func (parser *Parser) parseBodyFunctionBlock (
        parentIndent int,
        preExisting  *Block,
) (
        block *Block,
        err error,
) {
        if preExisting != nil {
                block = preExisting
        } else {
                block = &Block { variables: make(map[string] *Variable) }
        }
        
        block.where = parser.embedPosition()

        if (parser.line.Indent > 4) {
                parser.printWarning (
                        parser.token.Column,
                        "indentation level of",
                        parser.line.Indent,
                        "is difficult to read.",
                        "consider breaking up this function.")
        }

        for {
                if parser.line.Indent <= parentIndent {
                        break
                        
                } else if parser.line.Indent == parentIndent + 1 {
                        // we are parsing a statement
                        var statement *Statement
                        var worked bool
                        statement,
                        worked, err = parser.parseBodyFunctionStatement (
                                parentIndent + 1,
                                true, block)
                        if err != nil || !worked { return }

                        block.items = append (
                                block.items,
                                BlockOrStatement {
                                        statement: statement,
                                },
                        )
                        
                        if parser.endOfLine() {
                                parser.nextLine()
                        }
                        
                } else if parser.line.Indent == parentIndent + 2 {
                        // we are parsing a block
                        var childBlock *Block
                        childBlock, err = parser.parseBodyFunctionBlock (
                                parentIndent + 1, nil)
                        if err != nil { return }

                        block.items = append (block.items, BlockOrStatement {
                                block: childBlock,
                        })
                        
                } else {
                        parser.printError(0, errTooMuchIndent)
                        
                }

                if parser.endOfFile() { return }
        }

        return
}

/* parseBodyFunctionStatement parses a statement in a function body. This is
 * done recursively, and it may eat up more lines than one.
 */
func (parser *Parser) parseBodyFunctionStatement (
        parentIndent      int,
        isDirectlyInBlock bool,
        // specifically for defining variables in
        parent           *Block,
) (
        statement *Statement,
        worked bool,
        err error,
) {       
        statement = &Statement {
                where: parser.embedPosition(),
        }

        match := parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindSymbol)
        if !match {
                // we have no idea what the users intent with that was, so try
                // to move on to the next statement.
                parser.skipBodyFunctionStatement(parentIndent, false)
                return nil, false, nil
        }

        // if the first token found was a bracket, this statement is wrapped in
        // brackets and we have to do some things differently
        bracketed := parser.token.Kind == lexer.TokenKindLBracket
        if bracketed {
                // that wasn't the function name, so try to get the function
                // name again.
                parser.nextToken()
                match = parser.expect (
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                        lexer.TokenKindSymbol)
                if !match {
                        err = parser.skipBodyFunctionStatement (
                                parentIndent, bracketed)
                        return nil, false, err
                }
        }

        if parser.token.Kind == lexer.TokenKindString {
                // this statement calls a function of arbitrary name
                statement.external = true
                statement.externalCommand = parser.token.StringValue
                parser.nextToken()
        } else if parser.token.Kind == lexer.TokenKindSymbol {
                // this statement is an operator
                statement.command = Identifier {
                        trail: []string { parser.token.StringValue },
                }
                parser.nextToken()
        } else {
                // this statement calls a reachable function
                trail, worked, err := parser.parseIdentifier()
                        if err != nil { return nil, false, err }
                if !worked {
                        parser.skipBodyFunctionStatement(parentIndent, bracketed)
                        return nil, false, nil
                }

                statement.command = Identifier { trail: trail }
        }

        // get statement arguments
        complete := false
        for !complete {
                if !parser.expect (
                        lexer.TokenKindNone,
                        lexer.TokenKindDirection,
                        lexer.TokenKindLBracket,
                        lexer.TokenKindRBracket,
                        lexer.TokenKindLBrace,
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                        lexer.TokenKindRune,
                        lexer.TokenKindInteger,
                        lexer.TokenKindSignedInteger,
                        lexer.TokenKindFloat,
                ) {
                        err = parser.skipBodyFunctionStatement (
                                parentIndent, bracketed)
                        return nil, false, err
                }

                if (parser.token.Kind == lexer.TokenKindNone) {
                        // if we have brackets, we can continue to parse the
                        // statement on the next line. if we don't, we are done
                        // parsing this statement.
                        if bracketed {
                                parser.nextLine()
                        } else {
                                complete = true
                        }
                        continue
                } else if (parser.token.Kind == lexer.TokenKindRBracket) {
                        complete = true
                        parser.nextToken()
                        continue
                } else if (parser.token.Kind == lexer.TokenKindDirection) {
                        complete = true
                        continue
                }

                argument,
                worked, err := parser.parseBodyFunctionStatementArgument (
                        parentIndent, parent)
                if err != nil { return nil, false, err }
                if !worked { continue }

                statement.arguments = append (
                        statement.arguments,
                        argument)
        }

        // if we don't need to parse a return direction, stop
        if parser.token.Kind != lexer.TokenKindDirection || !isDirectlyInBlock {
                return statement, true, nil
        }

        // we need to parse a return direction
        parser.nextToken()
        for parser.expect (
                lexer.TokenKindName,
                lexer.TokenKindNone,
                lexer.TokenKindLBracket,
                lexer.TokenKindString,
        ) {
                if parser.token.Kind != lexer.TokenKindName { break }
                
                identifier,
                worked,
                err := parser.parseBodyFunctionIdentifierOrDeclaration(parent)
                if err != nil || !worked { return statement, false, err }

                statement.returnsTo = append(statement.returnsTo, identifier)
        }
        
        return statement, true, nil
}

func (parser *Parser) parseBodyFunctionStatementArgument (
        parentIndent int,
        parent *Block,
) (
        argument Argument,
        worked bool,
        err error,
) {
        switch parser.token.Kind {
        case lexer.TokenKindLBracket:
                childStatement,
                worked,
                err := parser.parseBodyFunctionStatement (
                        parentIndent,
                        false, parent)
                if err != nil { return argument, false, err }
                if !worked {
                        parser.nextToken()
                        return argument, false, nil
                }
                
                argument.kind = ArgumentKindStatement
                argument.statementValue = childStatement
                break
                
        case lexer.TokenKindLBrace:
                dereference,
                worked, err := parser.parseDereference(parentIndent, parent)
                if err != nil || !worked { return argument, false, err }

                argument.kind = ArgumentKindDereference
                argument.dereferenceValue = &dereference
                break
                                
        case lexer.TokenKindName:
                identifier,
                worked,
                err := parser.parseBodyFunctionIdentifierOrDeclaration(parent)
                if err != nil || !worked { return argument, false, err }

                argument.kind = ArgumentKindIdentifier
                argument.identifierValue = identifier
                break
                
        case lexer.TokenKindString:
                argument.kind = ArgumentKindString
                argument.stringValue = parser.token.StringValue
                parser.nextToken()
                break
                
        case lexer.TokenKindRune:
                argument.kind = ArgumentKindRune
                argument.runeValue = parser.token.Value.(rune)
                parser.nextToken()
                break
                
        case lexer.TokenKindInteger:
                argument.kind = ArgumentKindInteger
                argument.integerValue = parser.token.Value.(uint64)
                parser.nextToken()
                break
                
        case lexer.TokenKindSignedInteger:
                argument.kind = ArgumentKindSignedInteger
                argument.signedIntegerValue = parser.token.Value.(int64)
                parser.nextToken()
                break
                
        case lexer.TokenKindFloat:
                argument.kind = ArgumentKindFloat
                argument.floatValue = parser.token.Value.(float64)
                parser.nextToken()
                break
        }

        return argument, true, nil
}

/* parseBodyFunctionIdentifierOrDeclaration parses what may be either an
 * identifier or a variable declaration. It returns the identifier of what it
 * parsed, and if a variable was declared, it adds it to the parent block's data
 * section.
 */
func (parser *Parser) parseBodyFunctionIdentifierOrDeclaration (
        parent *Block,
) (
        identifier *Identifier,
        worked bool,
        err error,
) {
        trail, worked, err := parser.parseIdentifier()
        if err != nil { return nil, false, err }
        if !worked {
                parser.nextToken()
                return nil, false, nil
        }

        identifier = &Identifier {
                where: parser.embedPosition(),
                trail: trail,
        }
        if (parser.token.Kind != lexer.TokenKindColon) {
                return identifier, true, nil
        }

        if len(trail) != 1 {
                parser.printError (
                        parser.token.Column,
                        "cannot use member selection in declaration, name ",
                        "cannot have dots in it")
                return nil, false, nil
        }

        parser.nextToken()
        if !parser.expect (
                lexer.TokenKindLBrace,
                lexer.TokenKindName,
        ) { return nil, false, nil }

        what, worked, err := parser.parseType()
        if err != nil || !worked { return nil, false, err }

        name := trail[0]
        variable := &Variable {
                where: parser.embedPosition(),
                
                name: name,
                what: what,
        }

        // TODO: check all scopes above this
        if !parent.addVariable(variable) {
                parser.printError (
                        parser.token.Column,
                        "a variable with the name", name, "ia already defined",
                        "in this block")
                return nil, false, err
        }
        return identifier, true, nil
}

/* parseDereference parses a dereference of a value of the form {value N} where
 * N is an optional offset.
 */
func (parser *Parser) parseDereference (
        parentIndent int,
        parent *Block,
) (
        dereference Dereference,
        worked bool,
        err error,
) {
        if !parser.expect (lexer.TokenKindLBrace) {
                return dereference, false, nil
        }
        parser.nextToken()

        if (parser.token.Kind == lexer.TokenKindNone) {
                // if we are at the end of the line, just go on to the next one
                parser.nextLine()
        }
        
        if !parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindLBrace,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindInteger,
        ) {
                parser.skipBodyFunctionDereference()
                return dereference, false, nil
        }
        
        argument, worked, err := parser.parseBodyFunctionStatementArgument (
                parentIndent, parent)
        if err != nil || !worked {
                parser.skipBodyFunctionDereference()
                return dereference, false, err
        }

        dereference.dereferences = &argument

        if !parser.expect (
                lexer.TokenKindRBrace,
                lexer.TokenKindInteger,
        ) {
                parser.skipBodyFunctionDereference()
                return dereference, false, nil
        }
        
        // get the count, if there is one
        if parser.token.Kind == lexer.TokenKindInteger {
                dereference.offset = parser.token.Value.(uint64)
                parser.nextToken()
                if !parser.expect(lexer.TokenKindRBrace) {
                        parser.skipBodyFunctionDereference()
                        return dereference, false, nil
                }
        }
        
        parser.nextToken()

        return dereference, true, nil
}

/* skipBodyFunctionBlock skips to the end of the current block. This is done
 * soley based on indentation.
 */
func (parser *Parser) skipBodyFunctionBlock (parentIndent int) (err error) {
        for {
                parser.nextLine()
                if parser.endOfFile() { break }
                if parser.line.Indent <= parentIndent { break }
        }

        return
}

/* skipBodyFunctionStatement skips to the end of the current statement, or
 * indentation drop.
 */
func (parser *Parser) skipBodyFunctionStatement (
        parentIndent int,
        bracketed bool,
) (
        err error,
) {
        depth := 1
        
        if !bracketed {
                depth --
        }
        
        for {
                if parser.endOfFile() { return }
                if parser.endOfLine() { parser.nextLine() }

                if parser.token.Kind == lexer.TokenKindLBracket { depth ++ }
                if parser.token.Kind == lexer.TokenKindRBracket { depth -- }

                // if we drop out of the block or exit the statement, we are
                // done.
                if parser.line.Indent < parentIndent { break }
                if depth == 0 { break }
        
                parser.nextToken()
        }

        parser.nextLine()
        return
}

/* skipBodyFunctionDereference skips to the end of the current dereference.
 */
func (parser *Parser) skipBodyFunctionDereference () (err error) {
        depth := 1
        
        for {
                if parser.endOfFile() { return }
                if parser.endOfLine() { parser.nextLine() }

                if parser.token.Kind == lexer.TokenKindLBrace { depth ++ }
                if parser.token.Kind == lexer.TokenKindRBrace { depth -- }

                if depth == 0 { break }
        
                parser.nextToken()
        }

        parser.nextToken()
        return
}
