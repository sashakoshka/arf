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

// Length returns the length of the type if the type is a fixed length array.
// Otherwise, it just returns zero.
func (what Type) Length () (length uint64) {
	if what.kind == TypeKindArray {
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

// Values returns an iterator for the initialization values
func (values ObjectInitializationValues) Sections () (
	iterator types.Iterator[Argument],
) {
	iterator = types.NewIterator(values.attributes)
	return
}

// Length returns the amount of values
func (values ArrayInitializationValues) Length () (length int) {
	length = len(values.values)
	return
}

// Item returns the value at index
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
