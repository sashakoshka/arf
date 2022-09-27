package parser

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

// nameable allows a tree node to have a name.
type nameable struct {
	name string
}

// Name returns the name of the node.
func (node nameable) Name () (name string) {
	name = node.name
	return
}
// typeable allows a node to have a type.
type typeable struct {
	what Type
}

// Type returns the type of the node.
func (node typeable) Type () (what Type) {
	what = node.what
	return
}

// permissionable allows a node to have a permission.
type permissionable struct {
	permission types.Permission
}

// Permission returns the permision of the node.
func (node permissionable) Permission () (permission types.Permission) {
	permission = node.permission
	return
}

// valuable allows a node to have an argument value.
type valuable struct {
	value Argument
}

// Value returns the value argument of the node.
func (node valuable) Value () (value Argument) {
	value = node.value
	return
}

// multiValuable allows a node to have several argument values.
type multiValuable struct {
	values []Argument
}

// Value returns the value at index.
func (node multiValuable) Value (index int) (value Argument) {
	value = node.values[index]
	return
}

// Length returns the amount of values in the mode.
func (node multiValuable) Length () (length int) {
	length = len(node.values)
	return
}
