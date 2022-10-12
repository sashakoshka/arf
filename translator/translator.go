package translator

import "git.tebibyte.media/arf/arf/analyzer"

// Translate takes in a path to a module and an io.Writer, and outputs the
// corresponding C through the writer. The C code will import nothing and
// function as a standalone translation unit.
func Translate (modulePath string, output io.Writer) (err error) {
	// TODO
	return
}
