package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// TypeSection represents a type definition section.
type TypeSection struct {
	sectionBase
	what     Type
	complete bool
	argument Argument
	// TODO: do not add members from parent type. instead have a member
	// function to discern whether this type contains a particular member,
	// and have it recurse all the way up the family tree. it will be the
	// translator's job to worry about what members are placed where.
	members []ObjectMember
}

// ObjectMember is a member of an object type. 
type ObjectMember struct {
	name string
	
	// even if there is a private permission in another module, we still
	// need to include it in the semantic analysis because we need to know
	// what members objects have.
	permission types.Permission

	what Type
}

func (member ObjectMember) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		member.name, " ",
		member.permission.ToString(),
		"\n")
	output += member.what.ToString(indent + 1)
	return
}

// ToString returns all data stored within the type section, in string form.
func (section TypeSection) ToString (indent int) (output string) {
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

// analyzeTypeSection analyzes a type section.
func (analyzer AnalysisOperation) analyzeTypeSection () (
	section Section,
	err error,
) {
	outputSection := TypeSection { }
	outputSection.where = analyzer.currentPosition

	section = &outputSection
	analyzer.addSection(section)

	inputSection := analyzer.currentSection.(parser.TypeSection)
	if inputSection.Permission() == types.PermissionReadWrite {
		err = inputSection.NewError (
			"read-write (rw) permission not understood in this " +
			"context, try read-only (ro)",
			infoerr.ErrorKindError)
	}

	outputSection.permission = inputSection.Permission()

	outputSection.what, err = analyzer.analyzeType(inputSection.Type())
	if err != nil { return }

	if !inputSection.Argument().Nil() {
		outputSection.argument,
		err = analyzer.analyzeArgument(inputSection.Argument())
		if err != nil { return }

		// type check default value
		err = analyzer.typeCheck (
			inputSection.Argument().Location(),
			outputSection.argument,
			outputSection.what)
		if err != nil { return }
	}

	// TODO: analyze members

	outputSection.complete = true
	return
}
