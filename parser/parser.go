package parser

import "io"
import "os"
import "path/filepath"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/lexer"
import "git.tebibyte.media/arf/arf/infoerr"

// parsingOperation holds information about an ongoing parsing operation.
type parsingOperation struct {
	modulePath string
	token      lexer.Token
	tokens     []lexer.Token
	tokenIndex int
	skimming   bool

	tree SyntaxTree
}

// Fetch returns the parsed module located at the specified path, and returns an
// abstract syntax tree. If the module has not yet been parsed, it parses it
// first.
func Fetch (modulePath string, skim bool) (tree SyntaxTree, err error) {
	if modulePath[0] != '/' {
		panic("module path did not begin at filesystem root")
	}

	// try to hit cache
	cached, exists := cache[modulePath]
	if exists && !(!skim && cached.skimmed){
		tree = cached.tree
		return
	}

	// miss, so parse the module.
	parser := parsingOperation {
		modulePath: modulePath,
		skimming:   skim,
		tree: SyntaxTree {
			requires: make(map[string] string),
			sections: make(map[string] Section),			
		},
	}

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

	// cache tree
	cache[modulePath] = cacheItem {
		tree:    tree,
		skimmed: false,
	}
	
	return
}

// parse parses a file and adds it to the syntax tree.
func (parser *parsingOperation) parse (sourceFile *file.File) (err error) {
	var tokens []lexer.Token
	tokens, err = lexer.Tokenize(sourceFile)
	if err != nil { return }

	// reset the parser
	if len(tokens) == 0 { return }
	parser.tokens = tokens
	parser.token  = tokens[0]
	parser.tokenIndex = 0

	err = parser.parseMeta()
	if err != nil { return }

	err = parser.parseBody()
	if err != nil { return }

	return
}

// expect takes in a list of allowed token kinds, and returns an error if the
// current token isn't one of them. If the length of allowed is zero, this
// function will not return an error.
func (parser *parsingOperation) expect (allowed ...lexer.TokenKind) (err error) {
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

	err = infoerr.NewError (
		parser.token.Location(),
		message, infoerr.ErrorKindError)
	return
}

// nextToken is the same as expect, but it advances to the next token first.
func (parser *parsingOperation) nextToken (allowed ...lexer.TokenKind) (err error) {
	parser.tokenIndex ++
	if parser.tokenIndex >= len(parser.tokens) { return io.EOF }
	parser.token = parser.tokens[parser.tokenIndex]
	
	err = parser.expect(allowed...)
	return
}

// previousToken goes back one token. If the parser is already at the beginning,
// this does nothing.
func (parser *parsingOperation) previousToken () {
	parser.tokenIndex --
	if parser.tokenIndex < 0 { parser.tokenIndex = 0 }
	parser.token = parser.tokens[parser.tokenIndex]
	return
}

// skipIndentLevel advances the parser, ignoring every line with an indentation
// equal to or greater than the specified indent.
func (parser *parsingOperation) skipIndentLevel (indent int) (err error) {
	braceLevel   := 0
	parenLevel   := 0
	bracketLevel := 0

	for {
		if parser.token.Is(lexer.TokenKindNewline) {
			err = parser.nextToken()
			if err != nil { return }

			shouldBreak :=
				!parser.token.Is(lexer.TokenKindIndent) ||
				parser.token.Value().(int) < indent
			
			shouldBreak =
				shouldBreak      &&
				braceLevel   < 1 &&
				parenLevel   < 1 &&
				bracketLevel < 1

			if shouldBreak { return }
		}

		switch parser.token.Kind() {
		case lexer.TokenKindLBrace:   braceLevel ++
		case lexer.TokenKindRBrace:   braceLevel --
		case lexer.TokenKindLParen:   parenLevel ++
		case lexer.TokenKindRParen:   parenLevel --
		case lexer.TokenKindLBracket: bracketLevel ++
		case lexer.TokenKindRBracket: bracketLevel --
		}

		err = parser.nextToken()
		if err != nil { return }
	}
}

// skipWhitespace skips over newlines and indent tokens.
func (parser *parsingOperation) skipWhitespace () (err error) {
	for {
		isWhitespace :=
			parser.token.Is(lexer.TokenKindIndent) ||
			parser.token.Is(lexer.TokenKindNewline)

		if !isWhitespace {
			break
		}

		err = parser.nextToken()
		if err != nil { return }
	}

	return
}
