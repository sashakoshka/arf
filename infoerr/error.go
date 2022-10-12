package infoerr

import "os"
import "fmt"
import "git.tebibyte.media/arf/arf/file"

type ErrorKind int

const (
	ErrorKindError ErrorKind = iota
	ErrorKindWarn
)

type Error struct {
	file.Location 
	message string
	kind    ErrorKind
}

// NewError creates a new error at the specified location.
func NewError (
	location file.Location,
	message string,
	kind ErrorKind,
) (
	err Error,
) {
	if location.File() == nil {
		panic("cannot create new Error in a blank file")
	}
	
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
	if err.Width() > 0 {
		formattedMessage += fmt.Sprint (
			" \033[34m", err.Row() + 1,
			":", err.Column() + 1)
	}
	formattedMessage +=
		" \033[90min\033[0m " +
		err.File().Path() + "\n"
	
	if err.Width() > 0 {
		// print erroneous line
		line := err.File().GetLine(err.Row())
		formattedMessage += line + "\n"

		// position error marker
		var index int
		for index = 0; index < err.Column(); index ++ {
			if line[index] == '\t' {
				formattedMessage += "\t"
			} else {
				formattedMessage += " "
			}
		}
		
		// print an arrow with a tail spanning the width of the mistake
		for index < err.Column() + err.Width() - 1 {
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
