package arf

import "io"
import "github.com/sashakoshka/arf/parser"

func CompileModule (modulePath string, output io.Writer) (err error) {
	_, err = parser.Parse(modulePath)
	return
}
