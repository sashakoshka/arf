package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// EnumSection represents an enumerated type section.
type EnumSection struct {
	sectionBase
	what     Type
	members []EnumMember
}

// EnumMember is a member of an enumerated type. 
type EnumMember struct {
	locatable
	name string
	argument Argument
}

// ToString returns all data stored within the member, in string form.
func (member EnumMember) ToString (indent int) (output string) {
	output += doIndent(indent, "member ", member.name, "\n")
	if member.argument != nil {
		output += member.argument.ToString(indent + 1)
	}
	
	return
}

// ToString returns all data stored within the type section, in string form.
func (section EnumSection) ToString (indent int) (output string) {
	output += doIndent(indent, "enumSection ")
	output += section.permission.ToString() + " "
	output += section.where.ToString()
	output += "\n"
	output += section.what.ToString(indent + 1)
	for _, member := range section.members {
		output += member.ToString(indent + 1)
	}
	return
}

// analyzeEnumSection analyzes an enumerated type section.
func (analyzer analysisOperation) analyzeEnumSection () (
	section Section,
	err error,
) {
	outputSection := EnumSection { }
	outputSection.where = analyzer.currentPosition

	section = &outputSection
	analyzer.addSection(section)

	inputSection := analyzer.currentSection.(parser.EnumSection)
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

	// enum sections are only allowed to inherit from type sections
	_, inheritsFromTypeSection := outputSection.what.actual.(*TypeSection)
	if !inheritsFromTypeSection {
		err = inputSection.Type().NewError (
			"enum sections can only inherit from other type " +
			"sections.",
			infoerr.ErrorKindError)
		return
	}

	// analyze members
	for index := 0; index < inputSection.Length(); index ++ {
		inputMember := inputSection.Item(index)
		outputMember := EnumMember { }
		outputMember.location = inputMember.Location()
		outputMember.name     = inputMember.Name()

		if !inputMember.Argument().Nil() {
			outputMember.argument,
			err = analyzer.analyzeArgument(inputMember.Argument())
			if err != nil { return }
	
			// type check default value
			err = analyzer.typeCheck (
				outputMember.argument,
				outputSection.what)
			if err != nil { return }
		}
	}

	outputSection.complete = true
	return
}
