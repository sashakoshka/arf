package analyzer

// TypeSection represents a type definition section.
type TypeSection struct {
	sectionBase
	inherits Type
	complete bool
}

// Kind returns SectionKindType.
func (section TypeSection) Kind () (kind SectionKind) {
	kind = SectionKindType
	return
}

// ToString returns all data stored within the type section, in string form.
func (section TypeSection) ToString (indent int) (output string) {
	output += doIndent(indent, "typeSection ", section.where.ToString(), "\n")
	output += section.inherits.ToString(indent + 1)
	return
}
