package analyzer

func (section TypeSection) Kind () (kind SectionKind) {
	kind = SectionKindType
	return
}

func (section TypeSection) Name () (name string) {
	name = section.name
	return
}
