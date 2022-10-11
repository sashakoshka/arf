package analyzer

import "path/filepath"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

// locatable allows a tree node to have a location.
type locatable struct {
	location file.Location
}

// Location returns the location of the node.
func (node locatable) Location () (location file.Location) {
	location = node.location
	return
}

// NewError creates a new error at the node's location.
func (node locatable) NewError (
	message string,
	kind    infoerr.ErrorKind,
) (
	err error,
) {
	err = infoerr.NewError(node.location, message, kind)
	return
}

// sectionBase is a struct that all sections must embed.
type sectionBase struct {
	where      locator
	complete   bool
	permission types.Permission
	locatable
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
