package file

// Location represents a specific point in a file. It is used for error
// reporting.
type Location struct {
	file   *File
	row    int
	column int
	width  int
}

// NewError creates a new error at this location.
func (location Location) NewError (message string, kind ErrorKind) (err Error) {
	return NewError(location, message, kind)
}
