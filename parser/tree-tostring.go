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

	for _, name := range sortMapKeysAlphabetically(tree.requires) {
		require := tree.requires[name]
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
	if what.kind == TypeKindNil {
		output += "NIL-TYPE"
		return
	}

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

func (list List) ToString (indent int, breakline bool) (output string) {
	if !breakline { indent = 0 }
	output += doIndent(indent, "(")
	if breakline { output += "\n" }

	for index, argument := range list.arguments {
		if !breakline && index > 0 { output += " "}
		output += argument.ToString(indent, breakline)
	}

	output += doIndent(indent, ")")
	if breakline { output += "\n" }
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
	
	case ArgumentKindList:
		output += argument.value.(List).ToString(indent, breakLine)
	
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
	}

	return
}

func (section DataSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"data ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(), "\n")
	
	if section.argument.kind != ArgumentKindNil {
		output += section.argument.ToString(indent + 1, true)
	}
	
	if section.external {
		output += doIndent(indent + 1, "external\n")
	}
	
	return
}

func (member TypeSectionMember) ToString (indent int) (output string) {
	output += doIndent(indent, member.permission.ToString())
	output += " " + member.name

	if member.what.kind != TypeKindNil {
		output += ":" + member.what.ToString()
	}

	if member.argument.kind != ArgumentKindNil {
		output += " " + member.argument.ToString(indent, false)
	}

	if member.bitWidth > 0 {
		output += fmt.Sprint(" & ", member.bitWidth)
	}

	output += "\n"

	return
}

func (section TypeSection) ToString (indent int) (output string) {
	output += doIndent (
		indent,
		"type ",
		section.permission.ToString(), " ",
		section.name, ":",
		section.what.ToString(), "\n")
	
	if section.argument.kind != ArgumentKindNil {
		output += section.argument.ToString(indent + 1, true)
	}

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
		output += doIndent(indent + 1, "- ", member.name)
		if member.argument.kind != ArgumentKindNil {
			output += " " + member.argument.ToString(indent, false)
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

	if section.kind == FaceKindType {
		for _, name := range sortMapKeysAlphabetically(section.behaviors) {
			behavior := section.behaviors[name]
			output += behavior.ToString(indent + 1)
		}
	} else if section.kind == FaceKindFunc {
		for _, inputItem := range section.inputs {
			output += doIndent(indent + 1, "> ", inputItem.ToString(), "\n")
		}
		
		for _, outputItem := range section.outputs {
			output += doIndent(indent + 1, "< ", outputItem.ToString(), "\n")
		}
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

	output += "["

	switch phrase.kind {
	case PhraseKindCase:
        	output += ":"
        
	case PhraseKindAssign:
        	output += "="
	
	case PhraseKindOperator:
		
		switch phrase.operator {
	        case lexer.TokenKindColon:
	        	output += ":"
	        case lexer.TokenKindPlus:
	        	output += "+"
	        case lexer.TokenKindMinus:
	        	output += "-"
	        case lexer.TokenKindIncrement:
	        	output += "++"
	        case lexer.TokenKindDecrement:
	        	output += "--"
	        case lexer.TokenKindAsterisk:
	        	output += "*"
	        case lexer.TokenKindSlash:
	        	output += "/"
	        case lexer.TokenKindExclamation:
	        	output += "!"
	        case lexer.TokenKindPercent:
	        	output += "%"
	        case lexer.TokenKindPercentAssignment:
	        	output += "%="
	        case lexer.TokenKindTilde:
	        	output += "~"
	        case lexer.TokenKindTildeAssignment:
	        	output += "~="
	        case lexer.TokenKindAssignment:
	        	output += "="
	        case lexer.TokenKindEqualTo:
	        	output += "=="
	        case lexer.TokenKindNotEqualTo:
	        	output += "!="
	        case lexer.TokenKindLessThanEqualTo:
	        	output += "<="
	        case lexer.TokenKindLessThan:
	        	output += "<"
	        case lexer.TokenKindLShift:
	        	output += "<<"
	        case lexer.TokenKindLShiftAssignment:
	        	output += "<<="
	        case lexer.TokenKindGreaterThan:
	        	output += ">"
	        case lexer.TokenKindGreaterThanEqualTo:
	        	output += ">="
	        case lexer.TokenKindRShift:
	        	output += ">>"
	        case lexer.TokenKindRShiftAssignment:
	        	output += ">>="
	        case lexer.TokenKindBinaryOr:
	        	output += "|"
	        case lexer.TokenKindBinaryOrAssignment:
	        	output += "|="
	        case lexer.TokenKindLogicalOr:
	        	output += "||"
	        case lexer.TokenKindBinaryAnd:
	        	output += "&"
	        case lexer.TokenKindBinaryAndAssignment:
	        	output += "&="
	        case lexer.TokenKindLogicalAnd:
	        	output += "&&"
	        case lexer.TokenKindBinaryXor:
	        	output += "^"
	        case lexer.TokenKindBinaryXorAssignment:
	        	output += "^="
		}
		
	default:
		output += phrase.command.ToString(0, false)
	}
	
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

func (funcOutput FuncOutput) ToString (indent int) (output string) {
	output += doIndent(indent + 1)
	output += "< " + funcOutput.Declaration.ToString()
	if funcOutput.argument.kind != ArgumentKindNil {
		output += " " + funcOutput.argument.ToString(indent, false)
	}
	output += "\n"

	return
}

func (block Block) ToString (indent int) (output string) {
	for _, phrase := range block {
		output += phrase.ToString(indent, true)
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
		output += outputItem.ToString(indent)
	}

	output += doIndent(indent + 1, "---\n")

	if section.external {
		output += doIndent(indent + 1, "external\n")
	}
	
	output += section.root.ToString(indent + 1)
	return
}
