package parser

// LookupSection looks returns the section under the give name. If the section
// does not exist, nil is returned.
func (tree *SyntaxTree) LookupSection (name string) (section Section) {
	section = tree.sections[name]
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

// Length returns the length of the type if the type is an array. If the result
// is 0, this means the array has an undefined/variable length.
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
