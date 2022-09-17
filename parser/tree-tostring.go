package parser

import "fmt"
import "sort"
import "git.tebibyte.media/arf/arf/lexer"

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

func (tree SyntaxTree) ToString (indent int) (output string) {
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

	sectionKeys := sortMapKeysAlphabetically(tree.sections)
	for _, name := range sectionKeys {
		output += tree.sections[name].ToString(indent)
	}

	return
}

func (identifier Identifier) ToString () (output string) {
	for index, trailItem := range identifier.trail {
		if index > 0 {
			output += "."
		}

		output += trailItem
	}
	return
}

func (values ObjectDefaultValues) ToString (
	indent int,
	breakLine bool,
) (
	output string,
) {
	if !breakLine { indent = 0 }
	output += doIndent(indent, "(")
	if breakLine { output += "\n" }
	
	for _, name := range sortMapKeysAlphabetically(values) {
		value := values[name]
		
		output += doIndent(indent, "." + name + ":")
		
		isComplexDefaultValue :=
			value.kind == ArgumentKindObjectDefaultValues ||
			value.kind == ArgumentKindArrayDefaultValues
		
		if isComplexDefaultValue {
			if breakLine { output += "\n" }
			output += value.ToString(indent + 1, breakLine)
		} else {
			output += "<"
			output += value.ToString(indent + 1, false)
			output += ">"
		}
		if breakLine { output += "\n" }
	}
	output += doIndent(indent, ")")
	return
}

func (values ArrayDefaultValues) ToString (
	indent int,
	breakLine bool,
) (
	output string,
) {
	if !breakLine { indent = 0 }
	output += doIndent(indent, "<")
	if breakLine { output += "\n" }
	
	for _, value := range values {
		output += value.ToString(indent, breakLine)
	}
	
	output += doIndent(indent, ">")
	return
}

func (member TypeMember) ToString (indent int, breakLine bool) (output string) {
	output += doIndent(indent, ".")
	
	output += member.permission.ToString() + " "
	output += member.name + ":"
	output += member.what.ToString(indent + 1, breakLine)

	if member.bitWidth > 0 {
		output += fmt.Sprint(" & ", member.bitWidth)
	}

	if breakLine {
		output += "\n"
	}

	return
}

func (what Type) ToString (indent int, breakLine bool) (output string) {
	if what.kind == TypeKindBasic {
		output += what.name.ToString()
	} else {
		output += "{"
		output += what.points.ToString(indent, breakLine)

		if what.kind == TypeKindVariableArray {
			output += " .."
		}
		
		output += "}"
	}

	if what.length > 1 {
		output += fmt.Sprint(":", what.length)
	}
	
	if what.mutable {
		output += ":mut"
	}

	if what.members != nil {
		if breakLine {
			output += ":\n" + doIndent(indent, "(\n")
			for _, member := range what.members {
				output += member.ToString(indent, breakLine)
			}
			output += doIndent(indent, ")")
		} else {
			output += ":("
			for index, member := range what.members {
				if index > 0 { output += " " }
				output += member.ToString(indent, breakLine)
			}
			output += ")"
		}
	}

	defaultValueKind := what.defaultValue.kind
	if defaultValueKind != ArgumentKindNil {
		isComplexDefaultValue :=
			defaultValueKind == ArgumentKindObjectDefaultValues ||
			defaultValueKind == ArgumentKindArrayDefaultValues
		
		if isComplexDefaultValue {
			output += ":"
			if breakLine { output += "\n" }
			output += what.defaultValue.ToString(indent, breakLine)
		} else {
			output += ":<"
			output += what.defaultValue.ToString(indent, false)
			output += ">"
		}
	}
	return
}

func (declaration Declaration) ToString (indent int) (output string) {
	output += declaration.name + ":"
	output += declaration.what.ToString(indent, false)
	return
}

func (argument Argument) ToString (indent int, breakLine bool) (output string) {
	if !breakLine { indent = 0 }
	if argument.kind == ArgumentKindNil {
		output += "NIL-ARGUMENT"
		if breakLine { output += "\n" }
		return
	}

	switch argument.kind {
	case ArgumentKindPhrase:
		output += argument.value.(Phrase).ToString (	
				indent,
				breakLine)
	
	case ArgumentKindObjectDefaultValues:
		output += argument.value.(ObjectDefaultValues).
				ToString(indent, breakLine)
	
	case ArgumentKindArrayDefaultValues:
		output += argument.value.(ArrayDefaultValues).
				ToString(indent, breakLine)
	
	case ArgumentKindIdentifier:
		output += doIndent (
			indent,
			argument.value.(Identifier).ToString())
		if breakLine { output += "\n" }
	
	case ArgumentKindDeclaration:
		output += doIndent (
			indent,
			argument.value.(Declaration).ToString(indent))
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
		var stringValue string
		switch argument.value.(lexer.TokenKind) {
	        case lexer.TokenKindColon:
	        	stringValue = ":"
	        case lexer.TokenKindPlus:
	        	stringValue = "+"
	        case lexer.TokenKindMinus:
	        	stringValue = "-"
	        case lexer.TokenKindIncrement:
	        	stringValue = "++"
	        case lexer.TokenKindDecrement:
	        	stringValue = "--"
	        case lexer.TokenKindAsterisk:
	        	stringValue = "*"
	        case lexer.TokenKindSlash:
	        	stringValue = "/"
	        case lexer.TokenKindExclamation:
	        	stringValue = "!"
	        case lexer.TokenKindPercent:
	        	stringValue = "%"
	        case lexer.TokenKindPercentAssignment:
	        	stringValue = "%="
	        case lexer.TokenKindTilde:
	        	stringValue = "~"
	        case lexer.TokenKindTildeAssignment:
	        	stringValue = "~="
	        case lexer.TokenKindAssignment:
	        	stringValue = "="
	        case lexer.TokenKindEqualTo:
	        	stringValue = "=="
	        case lexer.TokenKindNotEqualTo:
	        	stringValue = "!="
	        case lexer.TokenKindLessThanEqualTo:
	        	stringValue = "<="
	        case lexer.TokenKindLessThan:
	        	stringValue = "<"
	        case lexer.TokenKindLShift:
	        	stringValue = "<<"
	        case lexer.TokenKindLShiftAssignment:
	        	stringValue = "<<="
	        case lexer.TokenKindGreaterThan:
	        	stringValue = ">"
	        case lexer.TokenKindGreaterThanEqualTo:
	        	stringValue = ">="
	        case lexer.TokenKindRShift:
	        	stringValue = ">>"
	        case lexer.TokenKindRShiftAssignment:
	        	stringValue = ">>="
	        case lexer.TokenKindBinaryOr:
	        	stringValue = "|"
	        case lexer.TokenKindBinaryOrAssignment:
	        	stringValue = "|="
	        case lexer.TokenKindLogicalOr:
	        	stringValue = "||"
	        case lexer.TokenKindBinaryAnd:
	        	stringValue = "&"
	        case lexer.TokenKindBinaryAndAssignment:
	        	stringValue = "&="
	        case lexer.TokenKindLogicalAnd:
	        	stringValue = "&&"
	        case lexer.TokenKindBinaryXor:
	        	stringValue = "^"
	        case lexer.TokenKindBinaryXorAssignment:
	        	stringValue = "^="
		}
		output += doIndent(indent, stringValue)
		if breakLine { output += "\n" }
	}

	return
}

func (section DataSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"data ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(indent + 1, true), "\n")
	
	if section.external {
		output += doIndent(indent + 1, "external\n")
	}
	
	return
}

func (section TypeSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"type ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(indent + 1, true), "\n")
	return
}


func (section EnumSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"enum ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(indent + 1, true), "\n")

	for _, member := range section.members {
		output += doIndent(indent + 1, "- ", member.name)
	
		isComplexInitialization :=
			member.value.kind == ArgumentKindObjectDefaultValues ||
			member.value.kind == ArgumentKindArrayDefaultValues
		
		if isComplexInitialization {
			output += "\n"
			output += member.value.ToString(indent + 2, true)
		} else if member.value.kind != ArgumentKindNil {
			output += " " + member.value.ToString(0, false)
		}
		output += "\n"
	}
	return
}

func (section FaceSection) ToString (indent int) (output string) {
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

func (behavior FaceBehavior) ToString (indent int) (output string) {
	output += doIndent(indent, behavior.name, "\n")
	
	for _, inputItem := range behavior.inputs {
		output += doIndent(indent + 1, "> ", inputItem.ToString(indent), "\n")
	}
	
	for _, outputItem := range behavior.outputs {
		output += doIndent(indent + 1, "< ", outputItem.ToString(indent), "\n")
	}

	return
}

func (phrase Phrase) ToString (indent int, ownLine bool) (output string) {
	if ownLine {
		output += doIndent(indent)
	}

	output += "[" + phrase.command.ToString(0, false)
	for _, argument := range phrase.arguments {
		output += " " + argument.ToString(0, false)
	}
	output += "]"

	if len(phrase.returnees) > 0 {
		output += " ->"
		for _, returnItem := range phrase.returnees {
			output += " " + returnItem.ToString(0, false)
		}
	}

	if ownLine {
		output += "\n"
		output += phrase.block.ToString(indent + 1)
	} else if len(phrase.block) > 0 {
		output += "NON-BLOCK-LEVEL-PHRASE-HAS-BLOCK"
	}
	return
}

func (block Block) ToString (indent int) (output string) {
	for _, phrase := range block {
		output += phrase.ToString(indent, true)
	}

	return
}

func (funcOutput FuncOutput) ToString (indent int) (output string) {
	output += doIndent (
		indent + 1,
		"< ", funcOutput.Declaration.ToString(indent), "\n")
	return
}

func (section FuncSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"func ",
		section.permission.ToString(), " ",
		section.name, "\n")
	
	if section.receiver != nil {
		output += doIndent (
			indent + 1,
			"@ ", section.receiver.ToString(indent), "\n")
	}
	
	for _, inputItem := range section.inputs {
		output += doIndent(indent + 1, "> ", inputItem.ToString(indent), "\n")
	}
	
	for _, outputItem := range section.outputs {
		output += outputItem.ToString(indent + 1)
	}

	output += doIndent(indent + 1, "---\n")

	if section.external {
		output += doIndent(indent + 1, "external\n")
	}
	
	output += section.root.ToString(indent + 1)
	return
}
