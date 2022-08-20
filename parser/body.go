package parser

import "git.tebibyte.media/sashakoshka/arf/lexer"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

// parse body parses the body of an arf file, after the metadata header.
func (parser *ParsingOperation) parseBody () (err error) {
	for {
		err = parser.expect(lexer.TokenKindName)
		if err != nil { return }

		sectionType := parser.token.Value().(string)
		switch sectionType {
		case "data":
			var section *DataSection
			section, err = parser.parseDataSection()
			if parser.tree.dataSections == nil {
				parser.tree.dataSections =
					make(map[string] *DataSection)
			}
			parser.tree.dataSections[section.name] = section
			if err != nil { return }
		case "type":
			var section *TypeSection
			section, err = parser.parseTypeSection()
			if parser.tree.typeSections == nil {
				parser.tree.typeSections =
					make(map[string] *TypeSection)
			}
			parser.tree.typeSections[section.name] = section
			if err != nil { return }
		case "face":
		case "enum":
		case "func":
		default:
			err = parser.token.NewError (
				"unknown section type \"" + sectionType + "\"",
				infoerr.ErrorKindError)
			return
		}
	}
}
