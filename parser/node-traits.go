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
	err infoerr.Error,
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
	value Argument
}

// Value returns the value argument of the node.
func (trait valuable) Value () (value Argument) {
	value = trait.value
	return
}
