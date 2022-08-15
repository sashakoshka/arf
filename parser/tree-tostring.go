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
		output += require.ToString(1)
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

func (attribute *ObjectAttribute) ToString (
	indent    int,
	breakLine bool,
) (
	output string,
) {
	if breakLine {
		output += doIndent(indent)
	}
	
	output += ", " +  attribute.name
	if breakLine {
		output += "\n" + attribute.value.ToString(indent + 1, true)
	} else {
		output += " " + attribute.value.ToString(0, false)
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

	switch argument.what {
	case ArgumentKindPhrase:
		output += doIndent (
			indent,
			argument.value.(*Phrase).ToString (	
				indent,
				breakLine))
	
	case ArgumentKindObjectAttribute:
		output += doIndent (
			indent,
			argument.value.(*ObjectAttribute).ToString (	
				indent,
				breakLine))
	
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

	// TODO: print out initialization values. if there is only one of them,
	// keep it on the same line. if there are more than one, give each its
	// own line.
	return
}
