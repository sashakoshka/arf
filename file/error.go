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
	message string
	kind    ErrorKind
}

// NewError creates a new error at the specified location.
func NewError (
	location Location,
	message string,
	kind ErrorKind,
) (
	err Error,
) {
	return Error {
		Location: location,
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
		line := err.Location.file.lines[err.Location.row]
		formattedMessage +=
			err.Location.file.lines[err.Location.row] + "\n"

		// position error marker
		var index int
		for index = 0; index < err.Location.column; index ++ {
			if line[index] == '\t' {
				formattedMessage += "\t"
			} else {
				formattedMessage += " "
			}
		}
		
		// print an arrow with a tail spanning the width of the mistake
		for err.width > 1 {
			if line[index] == '\t' {
				formattedMessage += "--------"
			} else {
				formattedMessage += "-"
			}
			index ++
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

// Message returns the error's message string
func (err Error) Message () (message string) {
	return err.message
}

// Kind returns what kind of error the error is.
func (err Error) Kind () (kind ErrorKind) {
	return err.kind
}
