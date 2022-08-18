package file

import "fmt"

// Location represents a specific point in a file. It is used for error
// reporting.
type Location struct {
	file   *File
	row    int
	column int
	width  int
}

// File returns the file the location is in
func (location Location) File () (file *File) {
	return location.file
}

// Row returns the row the location is positioned at in the file, starting at
// zero.
func (location Location) Row () (row int) {
	return location.row
}

// Column returns the column the location is positioned at in the file, starting
// at zero.
func (location Location) Column () (column int) {
	return location.column
}

// Width returns the amount of runes spanned by the location, starting at row
// and column.
func (location Location) Width () (width int) {
	return location.width
}

// Describe generates a description of the location for debug purposes
func (location Location) Describe () (description string) {
	return fmt.Sprint (
		"in ", location.file.Path(),
		" row ", location.row,
		" column ", location.column,
		" width ", location.width)
}
