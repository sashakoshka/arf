package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// FuncSection represents a type definition section.
type FuncSection struct {
	sectionBase
	external bool	
}

// ToString returns all data stored within the function section, in string form.
func (section FuncSection) ToString (indent int) (output string) {
	output += doIndent(indent, "funcSection ")
	output += section.permission.ToString() + " "
	output += section.where.ToString()
	output += "\n"
	
	// TODO: arguments
	// TODO: root block
	return
}

// analyzeFuncSection analyzes a function section.
func (analyzer *analysisOperation) analyzeFuncSection () (
	section Section,
	err error,
) {
	outputSection := FuncSection { }
	outputSection.where = analyzer.currentPosition

	section = &outputSection
	analyzer.addSection(section)

	inputSection := analyzer.currentSection.(parser.FuncSection)
	outputSection.location = analyzer.currentSection.Location()

	// TODO: do not do this if it is a method
	if inputSection.Permission() == types.PermissionReadWrite {
		err = inputSection.NewError (
			"read-write (rw) permission not understood in this " +
			"context, try read-only (ro)",
			infoerr.ErrorKindError)
		return
	}

	outputSection.permission = inputSection.Permission()

	// TODO: analyze inputs and outputs and reciever
	
	if inputSection.External() {
		outputSection.external = true
			
	} else {
		// TODO: analyze root block if not nil
	}
	
	return
}
