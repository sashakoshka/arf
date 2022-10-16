package analyzer

import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// EnumSection represents an enumerated type section.
type EnumSection struct {
	sectionBase
	what     Type
	members  []EnumMember
	argument Argument
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
	if section.argument != nil {
		output += section.argument.ToString(indent + 1)
	}
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

	// if the inherited type is a single number, we take note of that here
	// because it will allow us to do things like automatically fill in
	// member values if they are not specified.
	isNumeric :=
		outputSection.what.isNumeric() &&
		outputSection.what.isSingular()

	// enum sections are only allowed to inherit from type sections
	_, inheritsFromTypeSection := outputSection.what.actual.(*TypeSection)
	if !inheritsFromTypeSection {
		err = inputSection.Type().NewError (
			"enum sections can only inherit from other type " +
			"sections",
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

			// attempt to resolve the argument to a single constant
			// literal
			outputMember.argument, err =
				outputMember.argument.Resolve()
			if err != nil { return }
	
			// type check value
			err = analyzer.typeCheck (
				outputMember.argument,
				outputSection.what)
			if err != nil { return }
		} else if !isNumeric {
			// non-numeric enums must have filled in values
			err = inputMember.NewError (
				"member value must be specified manually for " +
				"non-numeric enums",
				infoerr.ErrorKindError)
			return
		}

		for _, compareMember := range outputSection.members {
			if compareMember.name == outputMember.name {
				err = inputMember.NewError (
					"enum member names must be unique",
					infoerr.ErrorKindError)
				return
			}

			if outputMember.argument  == nil { continue }
			if compareMember.argument == nil { continue }

			if compareMember.argument.Equals (
				outputMember.argument.Value(),
			) {
				err = inputMember.NewError (
					"enum member values must be unique",
					infoerr.ErrorKindError)
				return
			}
		}

		outputSection.members = append (
			outputSection.members,
			outputMember)
	}

	// fill in members that do not have values
	if isNumeric {
		for index, fillInMember := range outputSection.members {
			if fillInMember.argument != nil { continue }

			max := uint64(0)
			for _, compareMember := range outputSection.members {
			
				compareValue := compareMember.argument
				switch compareValue.(type) {
				case IntLiteral:
					number := uint64 (
						compareValue.(IntLiteral).value)
					if number > max {
						max = number
					}
				case UIntLiteral:
					number := uint64 (
						compareValue.(UIntLiteral).value)
					if number > max {
						max = number
					}
				case FloatLiteral:
					number := uint64 (
						compareValue.(FloatLiteral).value)
					if number > max {
						max = number
					}
				case nil:
					// do nothing
				default:
					panic (
						"invalid state: illegal " +
						"argument type while " +
						"attempting to fill in enum " +
						"member value for " +
						fillInMember.name + " in " +
						outputSection.location.Describe())
				}
			}
			
			// fill in
			fillInMember.argument = UIntLiteral {
				locatable: fillInMember.locatable,
				value: max + 1,
			}
			outputSection.members[index] = fillInMember
		}
	}

	if len(outputSection.members) < 1 {
		err = outputSection.NewError (
			"cannot create an enum with no members",
			infoerr.ErrorKindError)
		return
	}

	outputSection.argument = outputSection.members[0].argument

	outputSection.complete = true
	return
}
