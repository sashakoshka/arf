package parser

import "fmt"
import "sort"

func doIndent (indent int, input ...string) (output string) {
	for index := 0; index < indent; index ++ {
		output += "\t"
	}
	for _, inputSection := range input {
		output += inputSection
	}
	return
}

func sortMapKeysAlphabetically[KEY_TYPE any] (
	unsortedMap map[string] KEY_TYPE,
) (
	sortedKeys []string,
) {
	sortedKeys = make([]string, len(unsortedMap))
	index := 0
	for key, _ := range unsortedMap {
		sortedKeys[index] = key
		index ++
	}
	sort.Strings(sortedKeys)

	return
}

func (tree *SyntaxTree) ToString (indent int) (output string) {
	output += doIndent(indent, ":arf\n")

	if tree.author != "" {
		output += doIndent(indent, "author \"", tree.author, "\"\n")
	}

	if tree.license != "" {
		output += doIndent(indent, "license \"", tree.license, "\"\n")
	}

	for _, require := range tree.requires {
		output += doIndent(indent, "require \"", require, "\"\n")
	}
	
	output += doIndent(indent, "---\n")

	typeSectionKeys := sortMapKeysAlphabetically(tree.typeSections)
	for _, name := range typeSectionKeys {
		output += tree.typeSections[name].ToString(indent)
	}

	objtSectionKeys := sortMapKeysAlphabetically(tree.objtSections)
	for _, name := range objtSectionKeys {
		output += tree.objtSections[name].ToString(indent)
	}

	enumSectionKeys := sortMapKeysAlphabetically(tree.enumSections)
	for _, name := range enumSectionKeys {
		output += tree.enumSections[name].ToString(indent)
	}

	faceSectionKeys := sortMapKeysAlphabetically(tree.faceSections)
	for _, name := range faceSectionKeys {
		output += tree.faceSections[name].ToString(indent)
	}

	dataSectionKeys := sortMapKeysAlphabetically(tree.dataSections)
	for _, name := range dataSectionKeys {
		output += tree.dataSections[name].ToString(indent)
	}
	return
}

func (identifier *Identifier) ToString () (output string) {
	for index, trailItem := range identifier.trail {
		if index > 0 {
			output += "."
		}

		output += trailItem
	}
	return
}

func (what *Type) ToString () (output string) {
	if what.kind == TypeKindBasic {
		output += what.name.ToString()
	} else {
		output += "{"
		output += what.points.ToString()

		if what.kind == TypeKindArray {
			output += " "
			if what.length == 0 {
				output += ".."
			} else {
				output += fmt.Sprint(what.length)
			}
		}
		
		output += "}"
	}

	if what.mutable {
		output += ":mut"
	}
	
	return
}

func (declaration *Declaration) ToString () (output string) {
	output += declaration.name + ":"
	output += declaration.what.ToString()
	return
}

func (attributes *ObjectInitializationValues) ToString (
	indent int,
) (
	output string,
) {
	for _, name := range sortMapKeysAlphabetically(attributes.attributes) {
		value := attributes.attributes[name]
	
		output += doIndent(indent, ".", name)
		if value.kind == ArgumentKindObjectInitializationValues {
			output += "\n"
			output += value.ToString(indent + 1, true)
		} else {
			output += " " + value.ToString(0, false) + "\n"
		}
	}
	
	return
}

func (values *ArrayInitializationValues) ToString (
	indent int,
) (
	output string,
) {
	for _, value := range values.values {
		output += value.ToString(indent, true)
	}
	
	return
}

func (phrase *Phrase) ToString (indent int, breakLine bool) (output string) {
	if breakLine {
		output += doIndent (
			indent,
			"[", phrase.command.ToString(0, false))
		output += "\n"
		for _, argument := range phrase.arguments {
			output += doIndent (
				indent,
				argument.ToString(indent + 1, true))
		}
	} else {
		output += "[" + phrase.command.ToString(0, false)
		for _, argument := range phrase.arguments {
			output += " " + argument.ToString(0, false)
		}
	}
	
	output += "]"

	if len(phrase.returnsTo) > 0 {
		output += " ->"
		for _, returnItem := range phrase.returnsTo {
			output += " " + returnItem.ToString(0, false)
		}
	}

	if breakLine {
		output += "\n"
	}
	return
}

func (argument *Argument) ToString (indent int, breakLine bool) (output string) {
	if !breakLine { indent = 0 }
	if argument.kind == ArgumentKindNil {
		output += "NIL-ARGUMENT"
		if breakLine { output += "\n" }
		return
	}

	switch argument.kind {
	case ArgumentKindPhrase:
		output += argument.value.(*Phrase).ToString (	
				indent,
				breakLine)
	
	case ArgumentKindObjectInitializationValues:
		// this should only appear in contexts where breakLine is true
		output += argument.value.(*ObjectInitializationValues).
				ToString(indent)
	
	case ArgumentKindArrayInitializationValues:
		// this should only appear in contexts where breakLine is true
		output += argument.value.(*ArrayInitializationValues).
				ToString(indent)
	
	case ArgumentKindIdentifier:
		output += doIndent (
			indent,
			argument.value.(*Identifier).ToString())
		if breakLine { output += "\n" }
	
	case ArgumentKindDeclaration:
		output += doIndent (
			indent,
			argument.value.(*Declaration).ToString())
		if breakLine { output += "\n" }
	
	case ArgumentKindInt, ArgumentKindUInt, ArgumentKindFloat:
		output += doIndent(indent, fmt.Sprint(argument.value))
		if breakLine { output += "\n" }
	
	case ArgumentKindString:
		output += doIndent (
			indent,
			"\"" + argument.value.(string) + "\"")
		if breakLine { output += "\n" }
		
	case ArgumentKindRune:
		output += doIndent (
			indent,
			"'" + string(argument.value.(rune)) + "'")
		if breakLine { output += "\n" }
		
	case ArgumentKindOperator:
		// TODO
		// also when parsing this argument kind, don't do it in the
		// argument parsing function. do it specifically when parsing a
		// phrase command.
	}

	return
}

func (section *DataSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"data ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString())

	isComplexInitialization :=
		section.value.kind == ArgumentKindObjectInitializationValues ||
		section.value.kind == ArgumentKindArrayInitializationValues

	if section.value.value == nil {
		output += "\n"
	} else if isComplexInitialization {
		output += "\n"
		output += section.value.ToString(indent + 1, true)
	} else {
		output += " " + section.value.ToString(0, false)
		output += "\n"
	}
	return
}

func (section *TypeSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"type ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.inherits.ToString())

	isComplexInitialization :=
		section.defaultValue.kind == ArgumentKindObjectInitializationValues ||
		section.defaultValue.kind == ArgumentKindArrayInitializationValues

	if section.defaultValue.value == nil {
		output += "\n"
	} else if isComplexInitialization {
		output += "\n"
		output += section.defaultValue.ToString(indent + 1, true)
	} else {
		output += " " + section.defaultValue.ToString(0, false)
		output += "\n"
	}
	return
}

func (member ObjtMember) ToString (indent int) (output string) {
	output += doIndent(indent)
	
	output += member.permission.ToString() + " "
	output += member.name + ":"
	output += member.what.ToString()

	if member.bitWidth > 0 {
		output += fmt.Sprint(" & ", member.bitWidth)
	}
	
	isComplexInitialization :=
		member.defaultValue.kind == ArgumentKindObjectInitializationValues ||
		member.defaultValue.kind == ArgumentKindArrayInitializationValues
	
	if member.defaultValue.value == nil {
		output += "\n"
	} else if isComplexInitialization {
		output += "\n"
		output += member.defaultValue.ToString(indent + 1, true)
	} else {
		output += " " + member.defaultValue.ToString(0, false)
		output += "\n"
	}

	return
}

func (section *ObjtSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"objt ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.inherits.ToString(), "\n")

	for _, name := range sortMapKeysAlphabetically(section.members) {
		output += section.members[name].ToString(indent + 1)
	}
	return
}

func (section *EnumSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"enum ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(), "\n")

	for _, member := range section.members {
		output += doIndent(indent + 1, member.name)
	
		isComplexInitialization :=
			member.value.kind == ArgumentKindObjectInitializationValues ||
			member.value.kind == ArgumentKindArrayInitializationValues
		
		if member.value.value == nil {
			output += "\n"
		} else if isComplexInitialization {
			output += "\n"
			output += member.value.ToString(indent + 2, true)
		} else {
			output += " " + member.value.ToString(0, false)
			output += "\n"
		}
	}
	return
}

func (section *FaceSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"face ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.inherits.ToString(), "\n")

	for _, name := range sortMapKeysAlphabetically(section.behaviors) {
		behavior := section.behaviors[name]
		output += behavior.ToString(indent + 1)
	}
	return
}

func (behavior *FaceBehavior) ToString (indent int) (output string) {
	output += doIndent(indent, behavior.name, "\n")
	
	for _, inputItem := range behavior.inputs {
		output += doIndent(indent + 1, "> ", inputItem.ToString(), "\n")
	}
	
	for _, outputItem := range behavior.outputs {
		output += doIndent(indent + 1, "< ", outputItem.ToString(), "\n")
	}

	return
}
