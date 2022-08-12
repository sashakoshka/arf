package parser

import "os"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/lexer"

// ParsingOperation holds information about an ongoing parsing operation.
type ParsingOperation struct {
	modulePath string
	token      lexer.Token
	tokens     []lexer.Token
	tokenIndex int
}

// Parse reads the files located in the module specified by modulePath, and
// converts them into an abstract syntax tree.
func Parse (modulePath string) (tree *SyntaxTree, err error) {
	parser := ParsingOperation { modulePath: modulePath }
	tree, err = parser.parse()
	return
}

// parse runs the parsing operation.
func (parser *ParsingOperation) parse () (tree *SyntaxTree, err error) {
	if parser.modulePath[len(parser.modulePath) - 1] != '/' {
		parser.modulePath += "/"
	}

	var moduleFiles []os.DirEntry
	moduleFiles, err = os.ReadDir(parser.modulePath)
	if err != nil { return }

	for _, entry := range moduleFiles {
		if filepath.Ext(entry.Name()) != ".arf" || entry.IsDir() {
			continue
		}

		var sourceFile *file.File
		sourceFile, err = file.Open(parser.modulePath + entry.Name())
		if err != nil { return }

		var tokens []lexer.Token
		tokens, err = lexer.Tokenize(sourceFile)
		if err != nil { return }
		parser.tokens = append(parser.tokens, tokens...)
	}
	return
}
