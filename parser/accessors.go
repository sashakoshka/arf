package parser

import "git.tebibyte.media/arf/arf/types"

// LookupSection looks returns the section under the give name. If the section
// does not exist, nil is returned.
func (tree SyntaxTree) LookupSection (name string) (section Section) {
	section = tree.sections[name]
	return
}

// Sections returns an iterator for the tree's sections
func (tree SyntaxTree) Sections () (iterator types.Iterator[Section]) {
	iterator = types.NewIterator(tree.sections)
	return
}

// ResolveRequire returns the full path, from the filesystem root, of an import.
// This method will return false for exists if the module has not been
// imported.
func (tree SyntaxTree) ResolveRequire (name string) (path string, exists bool) {
	path, exists = tree.requires[name]
	return
}

// Kind returns the section's kind (SectionKindType).
func (section TypeSection) Kind () (kind SectionKind) {
	kind = SectionKindType
	return
}

// Kind returns the section's kind (SectionKindObjt).
func (section ObjtSection) Kind () (kind SectionKind) {
	kind = SectionKindObjt
	return
}

// Kind returns the section's kind (SectionKindEnum).
func (section EnumSection) Kind () (kind SectionKind) {
	kind = SectionKindEnum
	return
}

// Kind returns the section's kind (SectionKindFace).
func (section FaceSection) Kind () (kind SectionKind) {
	kind = SectionKindFace
	return
}

// Kind returns the section's kind (SectionKindData).
func (section DataSection) Kind () (kind SectionKind) {
	kind = SectionKindData
	return
}

// Kind returns the section's kind (SectionKindFunc).
func (section FuncSection) Kind () (kind SectionKind) {
	kind = SectionKindFunc
	return
}

// Length returns the amount of names in the identifier.
func (identifier Identifier) Length () (length int) {
	length = len(identifier.trail)
	return
}

// Item returns the name at the specified index.
func (identifier Identifier) Item (index int) (item string) {
	item = identifier.trail[index]
	return
}

// Kind returns the type's kind.
func (what Type) Kind () (kind TypeKind) {
	kind = what.kind
	return
}

// Mutable returns whether or not the type's data is mutable.
func (what Type) Mutable () (mutable bool) {
	mutable = what.mutable
	return
}

// Length returns the length of the type. If it is greater than 1, that means
// the type is a fixed length array.
func (what Type) Length () (length uint64) {
	length = 1
	if what.length > 1 {
		length = what.length
	}
	return
}

// Name returns the name of the type, if it is a basic type. Otherwise, it
// returns a zero value identifier.
func (what Type) Name () (name Identifier) {
	if what.kind == TypeKindBasic {
		name = what.name
	}
	return
}

// Points returns the type that this type points to, or is an array of. If the
// type is a basic type, this returns a zero value type.
func (what Type) Points () (points Type) {
	if what.kind != TypeKindBasic {
		points = *what.points
	}
	return
}

// Values returns an iterator for the initialization values.
func (values ObjectInitializationValues) Sections () (
	iterator types.Iterator[Argument],
) {
	iterator = types.NewIterator(values.attributes)
	return
}

// Length returns the amount of values.
func (values ArrayInitializationValues) Length () (length int) {
	length = len(values.values)
	return
}

// Item returns the value at index.
func (values ArrayInitializationValues) Value (index int) (value Argument) {
	value = values.values[index]
	return
}

// Kind returns what kind of argument it is.
func (argument Argument) Kind () (kind ArgumentKind) {
	kind = argument.kind
	return
}

// Value returns the underlying value of the argument. You can use Kind() to
// find out what to cast this to.
func (argument Argument) Value () (value any) {
	value = argument.value
	return
}

// BitWidth returns the bit width of the object member. If it is zero, it should
// be treated as unspecified.
func (member ObjtMember) BitWidth () (width uint64) {
	width = member.bitWidth
	return
}

// Length returns the amount of members in the section.
func (section ObjtSection) Length () (length int) {
	length = len(section.members)
	return
}

// Item returns the member at index.
func (section ObjtSection) Item (index int) (member ObjtMember) {
	member = section.members[index]
	return
}

// Length returns the amount of members in the section.
func (section EnumSection) Length () (length int) {
	length = len(section.members)
	return
}

// Item returns the member at index.
func (section EnumSection) Item (index int) (member EnumMember) {
	member = section.members[index]
	return
}

// InputsLength returns the amount of inputs in the behavior.
func (behavior FaceBehavior) IntputsLength () (length int) {
	length = len(behavior.inputs)
	return
}

// Input returns the input at index.
func (behavior FaceBehavior) Input (index int) (input Declaration) {
	input = behavior.inputs[index]
	return
}

// OutputsLength returns the amount of outputs in the behavior.
func (behavior FaceBehavior) OutputsLength () (length int) {
	length = len(behavior.outputs)
	return
}

// Output returns the output at index.
func (behavior FaceBehavior) Output (index int) (output Declaration) {
	output = behavior.outputs[index]
	return
}

// Behaviors returns an iterator for the interface's behaviors.
func (section FaceSection) Behaviors () (iterator types.Iterator[FaceBehavior]) {
	iterator = types.NewIterator(section.behaviors)
	return
}

// External returns whether or not the data section is external.
func (section DataSection) External () (external bool) {
	external = section.external
	return
}

// Kind returns what kind of phrase it is.
func (phrase Phrase) Kind () (kind PhraseKind) {
	kind = phrase.kind
	return
}

// ArgumentsLength returns the amount of arguments in the phrase.
func (phrase Phrase) ArgumentsLength () (length int) {
	length = len(phrase.arguments)
	return
}

// Argument returns the argument at index.
func (phrase Phrase) Argument (index int) (argument Argument) {
	argument = phrase.arguments[index]
	return
}

// ReturneesLength returns the amount of things the phrase returns to.
func (phrase Phrase) ReturneesLength () (length int) {
	length = len(phrase.returnees)
	return
}

// Returnee returns the returnee at index.
func (phrase Phrase) Returnee (index int) (returnee Argument) {
	returnee = phrase.returnees[index]
	return
}

// Block returns the block under the phrase, if it is a control flow statement.
func (phrase Phrase) Block () (block Block) {
	block = phrase.block
	return
}

// Receiver returns the method receiver, if there is one. Otherwise, it returns
// nil.
func (section FuncSection) Receiver () (receiver *Declaration) {
	receiver = section.receiver
	return
}

// InputsLength returns the number of inputs in the function.
func (section FuncSection) InputsLength () (length int) {
	length = len(section.inputs)
	return
}

// Input returns the input at index.
func (section FuncSection) Input (index int) (input Declaration) {
	input = section.inputs[index]
	return
}

// OutputsLength returns the number of outputs in the function.
func (section FuncSection) OutputsLength () (length int) {
	length = len(section.outputs)
	return
}

// Output returns the output at index.
func (section FuncSection) Output (index int) (output FuncOutput) {
	output = section.outputs[index]
	return
}

// Root returns the root block of the section.
func (section FuncSection) Root () (root Block) {
	root = section.root
	return
}

// External returns whether or not the function is an external function or not.
func (section FuncSection) External () (external bool) {
	external = section.external
	return
}
