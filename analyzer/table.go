package analyzer

// locator uniquely identifies a section in the section table.
type locator struct {
	modulePath string
	name string
}

func (where locator) ToString () (output string) {
	output += where.modulePath + "." + where.name
	return
}

// SectionTable stores a list of semantically analized sections from one module,
// and all sections that it requires from other modules.
type SectionTable map[locator] Section

// ToString returns the data stored in the table as a string.
func (table SectionTable) ToString (indent int) (output string) {
	for _, section := range table {
		output += section.ToString(indent)
	}

	return
} 

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
	where    locator
	complete bool
}

// Name returns the name of the section.
func (section sectionBase) Name () (name string) {
	name = section.where.name
	return
}

// ModulePath returns the full path of the module the section came from.
func (section sectionBase) ModulePath () (name string) {
	name = section.where.modulePath
	return
}

// TODO: ModuleName returns the name of the module where the section came from.

// Complete returns wether the section has been completed.
func (section sectionBase) Complete () (complete bool) {
	complete = section.complete
	return
}
