package analyzer

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// Argument represents a value that can be placed anywhere a value goes. This
// allows things like phrases being arguments to other phrases.
type Argument interface {
	// Phrase
	// List
	// Dereference
	// Variable
	// IntLiteral
	// UIntLiteral
	// FloatLiteral
	// StringLiteral

	What     () (what Type)
	Location () (location file.Location)
	NewError (message string, kind infoerr.ErrorKind) (err error)
	ToString (indent int) (output string)
	Equals   (value any) (equal bool)
	Value    () (value any)
	Resolve  () (constant Argument, err error)
	canBePassedAs (what Type) (allowed bool)
}

// phrase
// 	is what
// list
// 	is what
// dereference
// 	if length is greater than 1
// 		length is 1
// 		is what (ignore length)
// 	else
// 		is points of reduced of what
// variable
// 	is what
// int
// 	primitive is basic signed | float
// 	length is 1
// uint
// 	primitive is basic signed | unsigned | float
// 	length is 1
// float
// 	primitive is basic float
// 	length is 1
// string
// 	primitive is basic signed | unsigned | float
//	length is equal
//	or
//	reduced is variable array
//	reduced points to signed | unsigned | float
//	length is 1

// analyzeArgument analyzes an argument
func (analyzer analysisOperation) analyzeArgument (
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
		
	case parser.ArgumentKindList:
		// TODO
		
	case parser.ArgumentKindIdentifier:
		// TODO
		
	case parser.ArgumentKindDeclaration:
		// TODO
		
	case parser.ArgumentKindInt:
		outputArgument = IntLiteral {
			value: inputArgument.Value().(int64),
			locatable: locatable {
				location: inputArgument.Location(),
			},
		}
		
	case parser.ArgumentKindUInt:
		outputArgument = UIntLiteral {
			value: inputArgument.Value().(uint64),
			locatable: locatable {
				location: inputArgument.Location(),
			},
		}
		
	case parser.ArgumentKindFloat:
		outputArgument = FloatLiteral {
			value: inputArgument.Value().(float64),
			locatable: locatable {
				location: inputArgument.Location(),
			},
		}
		
	case parser.ArgumentKindString:
		outputArgument = StringLiteral {
			value: inputArgument.Value().(string),
			locatable: locatable {
				location: inputArgument.Location(),
			},
		}
	}
	return
}
