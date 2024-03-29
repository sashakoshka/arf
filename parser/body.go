package parser

import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parse body parses the body of an arf file, after the metadata header.
func (parser *parsingOperation) parseBody () (err error) {
	for {
		err = parser.expect(lexer.TokenKindName)
		if err != nil { return }
		sectionType := parser.token.Value().(string)

		switch sectionType {
		case "data":
			section, parseErr := parser.parseDataSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return parseErr }
			
		case "type":
			section, parseErr := parser.parseTypeSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return parseErr }
			
		case "face":
			section, parseErr := parser.parseFaceSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return parseErr }
			
		case "enum":
			section, parseErr := parser.parseEnumSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return parseErr }
			
		case "func":
			section, parseErr := parser.parseFuncSection()
			err = parser.tree.addSection(section)
			if err      != nil { return }
			if parseErr != nil { return parseErr }
			
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
	index := section.Name()

	funcSection, isFuncSection := section.(FuncSection)
	if isFuncSection {
		receiver := funcSection.receiver
		if receiver != nil {
			index = receiver.what.points.name.trail[0] + "_" + index
		}
	}

	_, exists := tree.sections[index]
	if exists {
		err = section.NewError (
			"cannot have multiple sections with the same name",
			infoerr.ErrorKindError)
		return
	}

	tree.sections[index] = section
	return
}
