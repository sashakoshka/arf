package parser

import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/infoerr"

// locatable allows a tree node to have a location.
type locatable struct {
	location file.Location
}

// Location returns the location of the node.
func (trait locatable) Location (location file.Location) {
	location = trait.location
	return
}

// setLocation sets the location of the node.
func (trait locatable) setLocation (location file.Location) {
	trait.location = location
}

// NewError creates a new error at the node's location.
func (trait locatable) NewError (
	message string,
	kind    infoerr.ErrorKind,
) (
	err infoerr.Error,
) {
	return infoerr.NewError(trait.location, message, kind)
}
