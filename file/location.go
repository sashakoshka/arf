package file

// Location represents a specific point in a file. It is used for error
// reporting.
type Location struct {
	file   *File
	row    int
	column int
}

// Error prints an error at this location.
func (location Location) Error (width int, message string) {
	location.file.Error(location.column, location.row, width, message)
}

// Warn prints a warning at this location.
func (location Location) Warn (width int, message string) {
	location.file.Warn(location.column, location.row, width, message)
}
