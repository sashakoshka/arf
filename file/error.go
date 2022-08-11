package file

import "os"
import "fmt"

type ErrorKind int

const (
	ErrorKindError ErrorKind = iota
	ErrorKindWarn
)

type Error struct {
	Location
	width   int
	message string
	kind    ErrorKind
}

// NewError creates a new error at the specified location.
func NewError (
	location Location,
	width int,
	message string,
	kind ErrorKind,
) (
	err *Error,
) {
	return &Error {
		Location: location,
		width:    width,
		message:  message,
		kind:     kind,
	}
}

// Error returns a formatted error message as a string.
func (err Error) Error () (formattedMessage string) {
	switch err.kind {
	case ErrorKindError:
		formattedMessage += "\033[31mERR\033[0m"
	case ErrorKindWarn:
		formattedMessage += "\033[33m!!!\033[0m"
	}

	// print information about the location of the mistake
	if err.width > 0 {
		formattedMessage += fmt.Sprint (
			" \033[34m", err.Location.row + 1,
			":", err.Location.column + 1)
	}
	formattedMessage +=
		" \033[90min\033[0m " +
		err.Location.file.path + "\n"
	
	if err.width > 0 {
		// print erroneous line
		formattedMessage +=
			err.Location.file.lines[err.Location.row] + "\n"

		// print an arrow with a tail spanning the width of the mistake
		columnCountdown := err.Location.column
		for columnCountdown > 1 {
			// TODO: for tabs, print out a teb instead.
			formattedMessage += " "
			columnCountdown --
		}
		for err.width > 1 {
			// TODO: for tabs, print out 8 of these instead.
			formattedMessage += "-"
		}
		formattedMessage += "^\n"
	}
	formattedMessage += err.message + "\n"

	return
}

// Print formats the error and prints it to stderr.
func (err Error) Print () {
	os.Stderr.Write([]byte(err.Error()))
}
