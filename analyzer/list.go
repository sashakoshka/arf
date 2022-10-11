package analyzer

// import "git.tebibyte.media/arf/arf/parser"
// import "git.tebibyte.media/arf/arf/infoerr"

type List struct {
	// TODO: length of what must be set to length of arguments
	what Type
	arguments []Argument
}

func (list List) ToString (indent int) (output string) {
	// TODO
	panic("TODO")
	return
} 

func (list List) canBePassedAs (what Type) (allowed bool) {
	// TODO
	panic("TODO")
	return
}
