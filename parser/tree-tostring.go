package parser

import "fmt"

func doIndent (indent int, input ...string) (output string) {
	for index := 0; index < indent; index ++ {
		output += "\t"
	}
	for _, inputSection := range input {
		output += inputSection
	}
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

	for _, require := range tree.dataSections {
		output += require.ToString(indent)
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
	for name, value := range attributes.attributes {
		output += doIndent(indent, ".", name, " ")
		if value.kind == ArgumentKindObjectInitializationValues {
			output += "\n"
			output += value.ToString(indent + 1, true)
		} else {
			output += value.ToString(0, false) + "\n"
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
	if argument.value == nil {
		output += "NIL ARGUMENT"
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
				ToString (indent)
	
	case ArgumentKindArrayInitializationValues:
		// this should only appear in contexts where breakLine is true
		output += argument.value.(*ArrayInitializationValues).
				ToString (indent)
	
	case ArgumentKindIdentifier:
		output += doIndent (
			indent,
			argument.value.(*Identifier).ToString())
	
	case ArgumentKindDeclaration:
		output += doIndent (
			indent,
			argument.value.(*Declaration).ToString())
	
	case ArgumentKindInt, ArgumentKindUInt, ArgumentKindFloat:
		output += doIndent(indent, fmt.Sprint(argument.value))
	
	case ArgumentKindString:
		output += doIndent (
			indent,
			"\"" + argument.value.(string) + "\"")
		
	case ArgumentKindRune:
		output += doIndent (
			indent,
			"'" + string(argument.value.(rune)) + "'")
		
	case ArgumentKindOperator:
		// TODO
		// also when parsing this argument kind, don't do it in the
		// argument parsing function. do it specifically when parsing a
		// phrase command.
	}

	if breakLine { output += "\n" }
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
