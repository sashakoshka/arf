package file

// Location represents a specific point in a file. It is used for error
// reporting.
type Location struct {
	file   *File
	row    int
	column int
}
