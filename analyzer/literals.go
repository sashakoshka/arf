package analyzer

import "fmt"

type IntLiteral    int64
type UIntLiteral   uint64
type FloatLiteral  float64
type StringLiteral string
type RuneLiteral   rune

// ToString outputs the data in the argument as a string.
func (literal IntLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg int ", literal, "\n"))
	return
}

// ToString outputs the data in the argument as a string.
func (literal UIntLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg uint ", literal, "\n"))
	return
}

// ToString outputs the data in the argument as a string.
func (literal FloatLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg float ", literal, "\n"))
	return
}

// ToString outputs the data in the argument as a string.
func (literal StringLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg string \"", literal, "\"\n"))
	return
}

// ToString outputs the data in the argument as a string.
func (literal RuneLiteral) ToString (indent int) (output string) {
	output += doIndent(indent, fmt.Sprint("arg rune '", literal, "'\n"))
	return
}

