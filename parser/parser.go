package parser

import "io"
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

	tree *SyntaxTree
}

// Parse reads the files located in the module specified by modulePath, and
// converts them into an abstract syntax tree.
func Parse (modulePath string) (tree *SyntaxTree, err error) {
	parser := ParsingOperation { modulePath: modulePath }

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

 		// parse the tokens into the module
		err  = parser.parse(sourceFile)
	}
	
	tree = parser.tree
	return
}

// parse parses a file and adds it to the syntax tree.
func (parser *ParsingOperation) parse (sourceFile *file.File) (err error) {
	var tokens []lexer.Token
	tokens, err = lexer.Tokenize(sourceFile)
	if err != nil { return }

	// reset the parser
	if parser.tree == nil {
		parser.tree = &SyntaxTree { }
	}
	if len(tokens) == 0 { return }
	parser.tokens = tokens
	parser.token  = tokens[0]
	parser.tokenIndex = 0

	err = parser.parseMeta()
	if err != nil { return }

	return
}

// expect takes in a list of allowed token kinds, and returns an error if the
// current token isn't one of them. If the length of allowed is zero, this
// function will not return an error.
func (parser *ParsingOperation) expect (allowed ...lexer.TokenKind) (err error) {
	if len(allowed) == 0 { return }

	for _, kind := range allowed {
		if parser.token.Is(kind) { return }
	}

	message :=
		"unexpected " + parser.token.Kind().Describe() +
		" token, expected "

	for index, allowedItem := range allowed {
		if index > 0 {
			if index == len(allowed) - 1 {
				message += " or "
			} else {
				message += ", " 
			}
		}
	
		message += allowedItem.Describe()
	}

	err = file.NewError (
		parser.token.Location(),
		message, file.ErrorKindError)
	return
}

// nextToken is the same as expect, but it advances to the next token first.
func (parser *ParsingOperation) nextToken (allowed ...lexer.TokenKind) (err error) {
	parser.tokenIndex ++
	if parser.tokenIndex >= len(parser.tokens) { return io.EOF }
	parser.token = parser.tokens[parser.tokenIndex]
	
	err = parser.expect(allowed...)
	return
}
