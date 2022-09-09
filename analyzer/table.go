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
	Kind () (kind SectionKind)
	Name () (name string)
	ToString () (output string)
}

// TypeKind represents what kind of type a type is.
type TypeKind int

const (
	// TypeKindBasic means it's a single value, or a fixed length array.
	TypeKindBasic TypeKind = iota

	// TypeKindPointer means it's a pointer
	TypeKindPointer

	// TypeKindVariableArray means it's an array of variable length.
	TypeKindVariableArray
)

// Type represents a description of a type. It must eventually point to a
// TypeSection.
type Type struct {
	actual Section
	points *Type

	mutable bool
	kind TypeKind
	length uint64
}

// TypeSection represents a type definition section.
type TypeSection struct {
	name     string
	inherits Type
}
