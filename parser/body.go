package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// TODO: parser must ensure that these names are unique

// parse body parses the body of an arf file, after the metadata header.
func (parser *ParsingOperation) parseBody () (err error) {
	for {
		err = parser.expect(lexer.TokenKindName)
		if err != nil { return }
		sectionType := parser.token.Value().(string)

		switch sectionType {
		case "data":
			section, parseErr := parser.parseDataSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		case "type":
			section, parseErr := parser.parseTypeSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		case "objt":
			section, parseErr := parser.parseObjtSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		case "face":
			section, parseErr := parser.parseFaceSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		case "enum":
			section, parseErr := parser.parseEnumSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		case "func":
			section, parseErr := parser.parseFuncSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return }
			
		default:
			err = parser.token.NewError (
				"unknown section type \"" + sectionType + "\"",
				infoerr.ErrorKindError)
			return
		}
	}
}

// addSection adds a section to the tree, ensuring it has a unique name within
// the module.
func (tree *SyntaxTree) addSection (section Section) (err error) {
	_, exists := tree.sections[section.Name()]
	if exists {
		err = section.NewError (
			"cannot have multiple sections with the same name",
			infoerr.ErrorKindError)
		return
	}

	tree.sections[section.Name()] = section
	return
}
