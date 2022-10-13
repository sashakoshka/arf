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
	locatable
	name string
	
	// even if there is a private permission in another module, we still
	// need to include it in the semantic analysis because we need to know
	// what members objects have.
	permission types.Permission

	what Type
	argument Argument
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

// Member returns the membrs ksdn ,mn ,mxc lkzxjclkjxzc l,mnzc .,zxmn.,zxmc
// IT RECURSES!
func (section TypeSection) Member (
	name string,
) (
	member ObjectMember,
	exists bool,
) {
	switch section.what.kind {
	case TypeKindBasic:
		for _, currentMember := range section.members {
			if currentMember.name == name {
				member = currentMember
				exists = true
				break
			}
		}
		
		if !exists {
			actual := section.what.actual
			if actual == nil { return }
			member, exists = actual.Member(name)
		}
		
	case TypeKindPointer:
		points := section.what.points
		if points == nil { return }
		member, exists = points.actual.Member(name)
	}
		
	return
}

// analyzeTypeSection analyzes a type section.
func (analyzer analysisOperation) analyzeTypeSection () (
	section Section,
	err error,
) {
	outputSection := TypeSection { }
	outputSection.where = analyzer.currentPosition

	section = &outputSection
	analyzer.addSection(section)

	inputSection := analyzer.currentSection.(parser.TypeSection)
	outputSection.location = analyzer.currentSection.Location()

	if inputSection.Permission() == types.PermissionReadWrite {
		err = inputSection.NewError (
			"read-write (rw) permission not understood in this " +
			"context, try read-only (ro)",
			infoerr.ErrorKindError)
		return
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
			outputSection.argument,
			outputSection.what)
		if err != nil { return }
	}

	// analyze members
	isObj := outputSection.what.underlyingPrimitive() == &PrimitiveObj
	if isObj {
		err = analyzer.analyzeObjectMembers (
			&outputSection,
			inputSection)
		if err != nil { return }
		
	} else if inputSection.MembersLength() > 0 {
		// if there are members, and the inherited type does not have
		// Obj as a primitive, throw an error.
		err = inputSection.Member(0).NewError (
			"members can only be defined on types descending " +
			"from Obj",
			infoerr.ErrorKindError)
		if err != nil { return }
	}

	outputSection.complete = true
	return
}

func (analyzer *analysisOperation) analyzeObjectMembers (
	into *TypeSection,
	from parser.TypeSection,
) (
	err error,
) {
	inheritedSection := into.what.actual
	for index := 0; index < from.MembersLength(); index ++ {
		inputMember := from.Member(index)

		outputMember := ObjectMember { }
		outputMember.location   = inputMember.Location()
		outputMember.name       = inputMember.Name()
		outputMember.permission = inputMember.Permission()
		
		inheritedMember, exists :=
			inheritedSection.Member(inputMember.Name())

		if exists {
			// modifying default value/permissions of an
			// inherited member

			outputMember.what = inheritedMember.what
			if !inputMember.Type().Nil() {
				err = inputMember.NewError (
					"cannot override type of " +
					"inherited member",
					infoerr.ErrorKindError)
				return
			}
			
			if outputMember.permission > inheritedMember.permission {
				err = inputMember.NewError (
					"cannot relax permission of " +
					"inherited member",
					infoerr.ErrorKindError)
				return
			}

			// apply default value
			if inputMember.Argument().Nil() {
				// if it is unspecified, inherit it
				outputMember.argument = inheritedMember.argument
			} else {
				outputMember.argument,
				err = analyzer.analyzeArgument(inputMember.Argument())
				if err != nil { return }

				// type check default value
				err = analyzer.typeCheck (
					outputMember.argument,
					outputMember.what)
				if err != nil { return }
			}
			
		} else {
			// defining a new member
			if inputMember.Type().Nil() {
				err = inputMember.NewError (
					"new members must be given a " +
					"type",
					infoerr.ErrorKindError)
				return
			}
			
			// apply default value
			if !inputMember.Argument().Nil() {
				outputMember.argument,
				err = analyzer.analyzeArgument(inputMember.Argument())
				if err != nil { return }

				// type check default value
				err = analyzer.typeCheck (
					outputMember.argument,
					outputMember.what)
				if err != nil { return }
			}
		}

		into.members = append (
			into.members,
			outputMember)
	}
	return
}
