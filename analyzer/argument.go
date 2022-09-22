package analyzer

import "git.tebibyte.media/arf/arf/parser"
// import "git.tebibyte.media/arf/arf/infoerr"

// Argument represents a value that can be placed anywhere a value goes. This
// allows things like phrases being arguments to other phrases.
type Argument interface {
	// Phrase
	// Dereference
	// Subscript
	// Object
	// Array
	// Variable
	// IntLiteral
	// UIntLiteral
	// FloatLiteral
	// StringLiteral
	// RuneLiteral

	ToString (indent int) (output string)
}

// analyzeArgument analyzes an argument
func (analyzer AnalysisOperation) analyzeArgument (
	inputArgument parser.Argument,
) (
	outputArgument Argument,
) {
	switch inputArgument.Kind() {
	case parser.ArgumentKindNil:
		
	case parser.ArgumentKindPhrase:
		
	case parser.ArgumentKindDereference:
		
	case parser.ArgumentKindSubscript:
		
	case parser.ArgumentKindObjectDefaultValues:
		
	case parser.ArgumentKindArrayDefaultValues:
		
	case parser.ArgumentKindIdentifier:
		
	case parser.ArgumentKindDeclaration:
		
	case parser.ArgumentKindInt:
		
	case parser.ArgumentKindUInt:
		
	case parser.ArgumentKindFloat:
		
	case parser.ArgumentKindString:
		
	case parser.ArgumentKindRune:
		
	case parser.ArgumentKindOperator:
		
	}
	return
}
