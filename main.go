package arf

import "io"
import "git.tebibyte.media/arf/arf/parser"

func CompileModule (modulePath string, output io.Writer) (err error) {
	_, err = parser.Parse(modulePath)
	return
}
