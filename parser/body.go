package parser

import "git.tebibyte.media/sashakoshka/arf/lexer"

// parse body parses the body of an arf file, after the metadata header.
func (parser *ParsingOperation) parseBody () (err error) {
	err = parser.nextToken(lexer.TokenKindName)
	if err != nil { return }

	switch parser.token.Value().(string) {
	case "data":
		var section *DataSection
		section, err = parser.parseDataSection()
		if err != nil { return }
		if parser.tree.dataSections == nil {
			parser.tree.dataSections = make(map[string] *DataSection)
		}
		parser.tree.dataSections[section.name] = section
	case "type":
	case "face":
	case "enum":
	case "func":
	}

	return
}
