package parser

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

// locatable allows a tree node to have a location.
type locatable struct {
	location file.Location
}

// Location returns the location of the node.
func (trait locatable) Location () (location file.Location) {
	location = trait.location
	return
}

// NewError creates a new error at the node's location.
func (trait locatable) NewError (
	message string,
	kind    infoerr.ErrorKind,
) (
	err error,
) {
	err = infoerr.NewError(trait.location, message, kind)
	return
}

// nameable allows a tree node to have a name.
type nameable struct {
	name string
}

// Name returns the name of the node.
func (trait nameable) Name () (name string) {
	name = trait.name
	return
}
// typeable allows a node to have a type.
type typeable struct {
	what Type
}

// Type returns the type of the node.
func (trait typeable) Type () (what Type) {
	what =  trait.what
	return
}

// permissionable allows a node to have a permission.
type permissionable struct {
	permission types.Permission
}

// Permission returns the permision of the node.
func (trait permissionable) Permission () (permission types.Permission) {
	permission = trait.permission
	return
}

// valuable allows a node to have an argument value.
type valuable struct {
	values []Argument
}

// Length returns the amount of default values in the node.
func (node valuable) Length () (length int) {
	length = len(node.values)
	return
} 

// Value returns the default value at the specified index.
func (node valuable) Value (index int) (value Argument) {
	value = node.values[index]
	return
}
