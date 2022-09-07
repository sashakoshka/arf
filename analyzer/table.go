package analyzer

// locator uniquely identifies a section in the section table.
type locator struct {
	modulePath string
	name string
}

// SectionTable stores a list of semantically analized sections from one module,
// and all sections that it requires from other modules.
type SectionTable map[locator] Section

// SectionKind differentiates Section interfaces.
type SectionKind int

const (
	SectionKindType = iota
	SectionKindObjt
	SectionKindEnum
	SectionKindFace
	SectionKindData
	SectionKindFunc
)

// Section is a semantically analyzed section.
type Section interface {
	Kind () (kind SectionKind)
	Name () (name string)
}
