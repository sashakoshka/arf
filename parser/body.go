package parser

import "github.com/sashakoshka/arf/lexer"

/* parseBody parses the body of an arf file. This contains sections, which have
 * code in them. Returns an error if the file cannot be parsed further.
 */
func (parser *Parser) parseBody (skim bool) (err error) {
        for !parser.endOfFile() {                
                if parser.line.Indent != 0 {
                        parser.printError(0, errBadIndent)
                        return nil
                }
                
                if !parser.expect(lexer.TokenKindName) { continue }

                switch parser.token.StringValue {
                case "data":
                        parser.nextToken()
                        section, err := parser.parseBodyData(skim, 0)
                        if err != nil { return err }
                        if section != nil {
                                err = parser.module.addData(section)
                                if err != nil { parser.printError(5, err) }
                        }
                        break
                case "type":
                        parser.nextToken()
                        section, err := parser.parseBodyTypedef(skim)
                        if err != nil { return err }
                        if section != nil {
                                err = parser.module.addTypedef(section)
                                if err != nil { parser.printError(5, err) }
                        }
                        break
                case "func":
                        parser.nextToken()
                        section, err := parser.parseBodyFunction(skim)
                        if err != nil { return err }
                        if section != nil {
                                err = parser.module.addFunction(section)
                                if err != nil { parser.printError(5, err) }
                        }
                        break
                default:
                        parser.printError (
                                0, "unknown section kind \"" +
                                parser.token.StringValue + "\"")
                        err = parser.skipBodySection()
                        if err != nil { return err }
                        break
                }
        }
        return
}

/* parseBodyData parses a data section.
 */
func (parser *Parser) parseBodyData (
        skim bool,
        parentIndent int,
) (
        section *Data,
        err error,
) {
        section = &Data {
                where: parser.embedPosition(),
        }

        if !parser.expect(lexer.TokenKindPermission) {
                return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)

        // if we are skimming and other modules don't have access to this, don't
        // parse it
        if (skim && section.modeExternal == ModeDeny) {
                section.external = true
                return nil, parser.skipBodySection()
        }
        
        worked := false
        section.name, section.what, worked, err = parser.parseDeclaration()
        if !worked {
                return nil, parser.skipBodySection()
        }
        
        // if we are skimming, don't parse the default values.
        if (skim) {
                section.external = true
                return section, parser.skipBodySection()
        }

        section.value, worked, err = parser.parseDefaultValues(parentIndent)
        if err != nil { return nil, err }
        if !worked { return nil, parser.skipBodySection() }

        return
}

/* parseBodyTypedef parses a type definition section.
 */
func (parser *Parser) parseBodyTypedef (
        skim bool,
) (
        section *Typedef,
        err error,
) {
        section = &Typedef {
                where: parser.embedPosition(),
        }

        if !parser.expect(lexer.TokenKindPermission) {
                 return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)
        
        // if we are skimming and other modules don't have access to this, don't
        // parse it
        if (skim && section.modeExternal == ModeDeny) {
                return nil, parser.skipBodySection()
        }

        worked := false
        section.name, section.inherits, worked, err = parser.parseDeclaration()
        if !worked {
                 return nil, parser.skipBodySection()
        }

        parser.nextToken()
        if !parser.expect() { return nil, parser.skipBodySection() }

        done := parser.nextLine()
        for {
                if done || parser.line.Indent == 0 { return }

                member, err := parser.parseBodyData(skim, 1)
                if err != nil { return nil, err }
                if member == nil { return nil, nil }

                section.members = append(section.members, member)
        }
}

/* skipBodySection ignores the rest of the current section of the body and moves
 * on to the next one.
 */
func (parser *Parser) skipBodySection () (err error ) {
        for {
                done := parser.nextLine()
                if done || parser.line.Indent == 0 { return }
        }
}


/* parseDefaultValues parses the default values of a variable.
 */
func (parser *Parser) parseDefaultValues (
        parentIndent int,
) (
        value []interface {},
        worked bool,
        err error,
) {
        for {
                if !parser.expect (
                        lexer.TokenKindNone,
                        lexer.TokenKindInteger,
                        lexer.TokenKindSignedInteger,
                        lexer.TokenKindFloat,
                        lexer.TokenKindString,
                        lexer.TokenKindRune,
                ) { return nil, false, nil }
                if parser.endOfLine() { break }
                value = append(value, parser.token.Value)
                parser.nextToken()
        }
        
        for {
                done := parser.nextLine()
                if done || parser.line.Indent <= parentIndent {
                        worked = true
                        return
                }

                for {
                        if !parser.expect (
                                lexer.TokenKindNone,
                                lexer.TokenKindInteger,
                                lexer.TokenKindSignedInteger,
                                lexer.TokenKindFloat,
                                lexer.TokenKindString,
                                lexer.TokenKindRune,
                        ) { return nil, false, nil }
                        if parser.endOfLine() { break }
                        value = append(value, parser.token.Value)
                        parser.nextToken()
                }
        }
}

/* parseDeclaration parses a variable declaration of the form name:Type or
 * name:{Type N}
 */
func (parser *Parser) parseDeclaration () (
        name string,
        what Type,
        worked bool,
        err error,
) {
        parser.nextToken()
        if !parser.expect(lexer.TokenKindName) { return }

        name = parser.token.StringValue

        parser.nextToken()
        if !parser.expect(lexer.TokenKindColon) { return }
        
        parser.nextToken()
        what, worked, err = parser.parseType()
        return
}

/* parseIdentifier parses an identifier of the form name.name.name
 */
func (parser *Parser) parseIdentifier () (
        trail []string,
        worked bool,
        err error,
) {
        for {
                if !parser.expect(lexer.TokenKindName) {
                        worked = false
                        return
                }
                
                trail = append(trail, parser.token.StringValue)
                
                parser.nextToken()

                if 
                        parser.endOfLine() ||
                        parser.token.Kind != lexer.TokenKindDot {
                
                        worked = true
                        return
                }
                parser.nextToken()
        }
}

/* parseType parses a type specifier that comes after a a colon.
 */
func (parser *Parser) parseType () (
        what   Type,
        worked bool,
        err    error,
) {
        if !parser.expect (
                lexer.TokenKindName,
                lexer.TokenKindLBrace,
        ) { return }
        
        // if the type is braced, we have a pointer
        if parser.token.Kind == lexer.TokenKindLBrace {
                parser.nextToken()

                // we must recurse to find what this type points to
                var typeThisPointsTo Type
                typeThisPointsTo, worked, err = parser.parseType()
                if !worked || err != nil { return }

                what.points = &typeThisPointsTo

                if !parser.expect (
                        lexer.TokenKindRBrace,
                        lexer.TokenKindInteger,
                ) {
                        return what, false, nil
                }
                
                // get the count, if there is one
                if parser.token.Kind == lexer.TokenKindInteger {
                        what.items = parser.token.Value.(uint64)
                        parser.nextToken()
                        if !parser.expect(lexer.TokenKindRBrace) {
                                return what, false, nil
                        }
                }
                
                parser.nextToken()

        // if the type is not braced, it is not a pointer
        } else {
                // get the identifier of this declaration's type
                what.name = Identifier {}
                
                what.name.trail, worked, err = parser.parseIdentifier()
                if !worked || err != nil { return }
        }
                
        // get an additional qualifier, if there is one
        if parser.token.Kind == lexer.TokenKindColon {
                parser.nextToken()
                if !parser.expect(lexer.TokenKindName) {
                        return what, false, nil
                }

                qualifier := parser.token.StringValue
                switch (qualifier) {
                case "mut":
                        what.mutable = true
                        break
                default:
                        parser.printError (
                                parser.token.Column,
                                "unknown type qualifier :" + qualifier)
                        break
                }
                
                parser.nextToken()
        }

        return what, true, nil
}

func decodePermission (value string) (internal Mode, external Mode) {
        if len(value) < 1 { return }
        switch value[0] {
                case 'n': internal = ModeDeny;  break
                case 'r': internal = ModeRead;  break
                case 'w': internal = ModeWrite; break
        }

        if len(value) < 2 { return }
        switch value[1] {
                case 'n': external = ModeDeny;  break
                case 'r': external = ModeRead;  break
                case 'w': external = ModeWrite; break
        }

        return
}
