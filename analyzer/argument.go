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
	canBePassedAs (what Type) (allowed bool)
}

// analyzeArgument analyzes an argument
func (analyzer AnalysisOperation) analyzeArgument (
	inputArgument parser.Argument,
) (
	outputArgument Argument,
	err error,
) {
	switch inputArgument.Kind() {
	case parser.ArgumentKindNil:
		panic("invalid state: attempt to analyze nil argument")
		
	case parser.ArgumentKindPhrase:
		// TODO
		
	case parser.ArgumentKindDereference:
		// TODO
		
	case parser.ArgumentKindSubscript:
		// TODO
		
	case parser.ArgumentKindList:
		// TODO
		
	case parser.ArgumentKindIdentifier:
		// TODO
		
	case parser.ArgumentKindDeclaration:
		// TODO
		
	case parser.ArgumentKindInt:
		outputArgument = IntLiteral(inputArgument.Value().(int64))
		
	case parser.ArgumentKindUInt:
		outputArgument = UIntLiteral(inputArgument.Value().(uint64))
		
	case parser.ArgumentKindFloat:
		outputArgument = FloatLiteral(inputArgument.Value().(float64))
		
	case parser.ArgumentKindString:
		outputArgument = StringLiteral(inputArgument.Value().(string))
		
	case parser.ArgumentKindRune:
		outputArgument = RuneLiteral(inputArgument.Value().(rune))
	}
	return
}
