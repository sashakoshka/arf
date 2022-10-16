package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// DataSection represents a global variable section.
type DataSection struct {
	sectionBase
	what     Type
	argument Argument
}

// ToString returns all data stored within the data section, in string form.
func (section DataSection) ToString (indent int) (output string) {
	output += doIndent(indent, "typeSection ")
	output += section.permission.ToString() + " "
	output += section.where.ToString()
	output += "\n"
	output += section.what.ToString(indent + 1)
	if section.argument != nil {
		output += section.argument.ToString(indent + 1)
	}
	return
}

// analyzeDataSection analyzes a data section.
func (analyzer analysisOperation) analyzeDataSection () (
	section Section,
	err error,
) {
	outputSection := DataSection { }
	outputSection.where = analyzer.currentPosition

	section = &outputSection
	analyzer.addSection(section)

	inputSection := analyzer.currentSection.(parser.DataSection)
	outputSection.location = analyzer.currentSection.Location()

	if inputSection.Permission() == types.PermissionReadWrite {
		err = inputSection.NewError (
			"read-write (rw) permission not understood in this " +
			"context, try read-only (ro)",
			infoerr.ErrorKindError)
		return
	}

	outputSection.permission = inputSection.Permission()

	// get inherited type
	outputSection.what, err = analyzer.analyzeType(inputSection.Type())
	if err != nil { return }

	// data sections are only allowed to inherit type, enum, and face sections
	_, inheritsFromTypeSection := outputSection.what.actual.(*TypeSection)
	_, inheritsFromEnumSection := outputSection.what.actual.(*EnumSection)
	// _, inheritsFromFaceSection := outputSection.what.actual.(*FaceSection)
	if !inheritsFromTypeSection && inheritsFromEnumSection {
		err = inputSection.Type().NewError (
			"type sections can only inherit from type, enum, and " +
			"face sections",
			infoerr.ErrorKindError)
		return
	}

	if !inputSection.Argument().Nil() {
		outputSection.argument,
		err = analyzer.analyzeArgument(inputSection.Argument())
		if err != nil { return }

		// type check default value
		err = analyzer.typeCheck (
			outputSection.argument,
			outputSection.what)
		if err != nil { return }
	}
	
	outputSection.complete = true
	return
}
