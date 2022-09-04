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
