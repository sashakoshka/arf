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
	SectionKindType SectionKind = iota
	SectionKindObjt
	SectionKindEnum
	SectionKindFace
	SectionKindData
	SectionKindFunc
)

// Section is a semantically analyzed section.
type Section interface {
	Kind     () (kind SectionKind)
	Name     () (name string)
	ToString (indent int) (output string)
	Complete () (complete bool)
}

// sectionBase is a struct that all sections must embed.
type sectionBase struct {
	name     string
	complete bool
}

// Name returns the name of the section.
func (section sectionBase) Name () (name string) {
	name = section.name
	return
}

// Complete returns wether the section has been completed.
func (section sectionBase) Complete () (complete bool) {
	complete = section.complete
	return
}
