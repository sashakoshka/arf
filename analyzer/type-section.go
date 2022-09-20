package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// TypeSection represents a type definition section.
type TypeSection struct {
	sectionBase
	what     Type
	complete bool
}

// Kind returns SectionKindType.
func (section TypeSection) Kind () (kind SectionKind) {
	kind = SectionKindType
	return
}

// ToString returns all data stored within the type section, in string form.
func (section TypeSection) ToString (indent int) (output string) {
	output += doIndent(indent, "typeSection ", section.where.ToString(), "\n")
	output += section.what.ToString(indent + 1)
	return
}

// analyzeTypeSection analyzes a type section.
func (analyzer AnalysisOperation) analyzeTypeSection () (
	section Section,
	err error,
) {
	inputSection := analyzer.currentSection.(parser.TypeSection)
	if inputSection.Permission() == types.PermissionReadWrite {
		err = inputSection.NewError (
			"rw permission not understood in this context, try ro",
			infoerr.ErrorKindError)
	}

	outputSection := TypeSection { }
	outputSection.where = analyzer.currentPosition

	outputSection.what, err = analyzer.analyzeType(inputSection.Type())
	if err != nil { return }

	outputSection.complete = true
	return
}
