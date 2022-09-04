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
		case "objt":
			var section *ObjtSection
			section, err = parser.parseObjtSection()
			if parser.tree.objtSections == nil {
				parser.tree.objtSections =
					make(map[string] *ObjtSection)
			}
			parser.tree.objtSections[section.name] = section
			if err != nil { return }
		case "face":
			var section *FaceSection
			section, err = parser.parseFaceSection()
			if parser.tree.faceSections == nil {
				parser.tree.faceSections =
					make(map[string] *FaceSection)
			}
			parser.tree.faceSections[section.name] = section
			if err != nil { return }
		case "enum":
			var section *EnumSection
			section, err = parser.parseEnumSection()
			if parser.tree.enumSections == nil {
				parser.tree.enumSections =
					make(map[string] *EnumSection)
			}
			parser.tree.enumSections[section.name] = section
			if err != nil { return }
		case "func":
			var section *FuncSection
			section, err = parser.parseFuncSection()
			if parser.tree.funcSections == nil {
				parser.tree.funcSections =
					make(map[string] *FuncSection)
			}
			parser.tree.funcSections[section.name] = section
			if err != nil { return }
		default:
			err = parser.token.NewError (
				"unknown section type \"" + sectionType + "\"",
				infoerr.ErrorKindError)
			return
		}
	}
}
