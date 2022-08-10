package arf

import "os"
import "io"
import "path/filepath"
// import "github.com/sashakoshka/arf/lexer"

func CompileModule (modulePath string, output io.Writer) (err error) {
	moduleFiles, err := os.ReadDir(modulePath)
	if err != nil { return err }

	// var moduleTokens []lexer.Token
	for _, entry := range moduleFiles {
		if filepath.Ext(entry.Name()) != ".arf" || entry.IsDir() {
			continue
		}

		// tokens, err := lexer.Tokenize()
		// if err != nil { return err }
	}

	return
}
