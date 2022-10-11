package analyzer

import "path/filepath"
import "git.tebibyte.media/arf/arf/types"

// sectionBase is a struct that all sections must embed.
type sectionBase struct {
	where      locator
	complete   bool
	permission types.Permission
}

// Name returns the name of the section.
func (section sectionBase) Name () (name string) {
	name = section.where.name
	return
}

// ModulePath returns the full path of the module the section came from.
func (section sectionBase) ModulePath () (path string) {
	path = section.where.modulePath
	return
}

// ModuleName returns the name of the module where the section came from.
func (section sectionBase) ModuleName () (name string) {
	name = filepath.Base(section.where.modulePath)
	return
}

// Complete returns wether the section has been completed.
func (section sectionBase) Complete () (complete bool) {
	complete = section.complete
	return
}

// Permission returns the permission of the section.
func (section sectionBase) Permission () (permission types.Permission) {
	permission = section.permission
	return
}

// locator returns the module path and name of the section.
func (section sectionBase) locator () (where locator) {
	where = section.where
	return
}
