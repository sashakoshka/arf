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

func (what Type) ToString () (output string) {
	if what.kind == TypeKindBasic {
		output += what.name.ToString()
	} else {
		output += "{"
		output += what.points.ToString()

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
	
	return
}

func (declaration Declaration) ToString () (output string) {
	output += declaration.name + ":"
	output += declaration.what.ToString()
	return
}

func (attributes ObjectInitializationValues) ToString (
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

func (values ArrayInitializationValues) ToString (
	indent int,
) (
	output string,
) {
	for _, value := range values.values {
		output += value.ToString(indent, true)
	}
	
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
	
	case ArgumentKindObjectInitializationValues:
		// this should only appear in contexts where breakLine is true
		output += argument.value.(ObjectInitializationValues).
				ToString(indent)
	
	case ArgumentKindArrayInitializationValues:
		// this should only appear in contexts where breakLine is true
		output += argument.value.(ArrayInitializationValues).
				ToString(indent)
	
	case ArgumentKindIdentifier:
		output += doIndent (
			indent,
			argument.value.(Identifier).ToString())
		if breakLine { output += "\n" }
	
	case ArgumentKindDeclaration:
		output += doIndent (
			indent,
			argument.value.(Declaration).ToString())
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

func (section TypeSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"type ",
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

func (member ObjtMember) ToString (indent int) (output string) {
	output += doIndent(indent)
	
	output += member.permission.ToString() + " "
	output += member.name + ":"
	output += member.what.ToString()

	if member.bitWidth > 0 {
		output += fmt.Sprint(" & ", member.bitWidth)
	}
	
	isComplexInitialization :=
		member.value.kind == ArgumentKindObjectInitializationValues ||
		member.value.kind == ArgumentKindArrayInitializationValues
	
	if member.value.value == nil {
		output += "\n"
	} else if isComplexInitialization {
		output += "\n"
		output += member.value.ToString(indent + 1, true)
	} else {
		output += " " + member.value.ToString(0, false)
		output += "\n"
	}

	return
}

func (section ObjtSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"objt ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.inherits.ToString(), "\n")

	for _, member := range section.members {
		output += member.ToString(indent + 1)
	}
	return
}

func (section EnumSection) ToString (indent int) (output string) {
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
		output += doIndent(indent + 1, "> ", inputItem.ToString(), "\n")
	}
	
	for _, outputItem := range behavior.outputs {
		output += doIndent(indent + 1, "< ", outputItem.ToString(), "\n")
	}

	return
}

func (phrase Phrase) ToString (indent int, ownLine bool) (output string) {
	if ownLine {
		output += doIndent(indent)
	}

	var initializationValues Argument
	
	output += "[" + phrase.command.ToString(0, false)
	for _, argument := range phrase.arguments {
		isInitializationValue :=
			argument.kind == ArgumentKindObjectInitializationValues ||
			argument.kind == ArgumentKindArrayInitializationValues
		if isInitializationValue {
			initializationValues = argument
		} else {
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

	if ownLine {
		output += "\n"
		if initializationValues.kind != ArgumentKindNil {
			output += initializationValues.ToString(indent + 1, true)
		}
		output += phrase.block.ToString(indent + 1)
	}
	return
}

func (block Block) ToString (indent int) (output string) {
	for _, phrase := range block {
		output += phrase.ToString(indent, true)
	}

	return
}

func (funcOutput FuncOutput) ToString () (output string) {
	output += funcOutput.Declaration.ToString()
	if funcOutput.value.kind != ArgumentKindNil {
		output += " " + funcOutput.value.ToString(0, false)
	}
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
			"@ ", section.receiver.ToString(), "\n")
	}
	
	for _, inputItem := range section.inputs {
		output += doIndent(indent + 1, "> ", inputItem.ToString(), "\n")
	}
	
	for _, outputItem := range section.outputs {
		output += doIndent(indent + 1, "< ", outputItem.ToString(), "\n")
	}

	output += doIndent(indent + 1, "---\n")

	if section.external {
		output += doIndent(indent + 1, "external\n")
	}
	
	output += section.root.ToString(indent + 1)
	return
}
