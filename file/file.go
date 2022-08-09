package file

import "os"
import "log"
import "bufio"

var logger = log.New(os.Stderr, "", 0)

// File represents a read only file that can print out formatted errors.
type File struct {
	path          string
	file          *os.File
	reader        *bufio.Reader
	currentLine   int
	currentColumn int
	lines         []string
}

// Open opens the file specified by path and returns a new File struct.
func Open (path string) (file *File, err error) {
	file = &File {
		path:  path,
		lines: []string { "" },
	}

	file.file, err = os.OpenFile(path, os.O_RDONLY, 0660)
	if err != nil { return }

	file.reader = bufio.NewReader(file.file)
	return
}

// Stat returns the FileInfo structure describing file. If there is an error, it
// will be of type *PathError. 
func (file *File) Stat () (fileInfo os.FileInfo, err error) {
	return file.file.Stat()
}

// Read reads up to len(bytes) bytes from the File and stores them in bytes. It
// returns the number of bytes read and any error encountered. At end of file,
// Read returns 0, io.EOF. This method stores all data read into the file struct
// so it can be used later for error reporting.
func (file *File) Read (bytes []byte) (amountRead int, err error) {
	amountRead, err = file.reader.Read(bytes)

	// store the character in the file
	for _, char := range bytes {
		if char == '\n' {
			file.lines = append(file.lines, "")
			file.currentLine ++
			file.currentColumn = 0
		} else {
			file.lines[file.currentLine] += string(char)
			file.currentColumn ++
		}
	}
	
	return
}

// ReadRune reads a single UTF-8 encoded Unicode character and returns the rune
// and its size in bytes. If the encoded rune is invalid, it consumes one byte
// and returns unicode.ReplacementChar (U+FFFD) with a size of 1.
func (file *File) ReadRune () (char rune, size int, err error) {
	char, size, err = file.reader.ReadRune()

	if char == '\n' {
		file.lines = append(file.lines, "")
		file.currentLine ++
		file.currentColumn = 0
	} else {
		file.lines[file.currentLine] += string(char)
		file.currentColumn ++
	}
	return
}

// ReadString reads until the first occurrence of delimiter in the input,
// returning a string containing the data up to and including the delimiter. If
// ReadString encounters an error before finding a delimiter, it returns the
// data read before the error and the error itself (often io.EOF). ReadString
// returns err != nil if and only if the returned data does not end in delim.
func (file *File) ReadString (delimiter byte) (read string, err error) {
	read, err = file.reader.ReadString(delimiter)

	// store the character in the file
	for _, char := range read {
		if char == '\n' {
			file.lines = append(file.lines, "")
			file.currentLine ++
			file.currentColumn = 0
		} else {
			file.lines[file.currentLine] += string(char)
			file.currentColumn ++
		}
	}
	
	return
}

// Close closes the file. After the file is closed, data that has been read will
// still be retained, and errors can be reported.
func (file *File) Close () {
	file.file.Close()
}

// mistake prints an informational message about a mistake that the user made.
func (file *File) mistake (column, row, width int, message, symbol string) {
	logger.Print(symbol)

	// print information about the location of the mistake
	if width > 0 {
		logger.Print(" ", row + 1, ":", column + 1)
	}
	logger.Println(" in", file.path)
	
	if width > 0 {
		// print erroneous line
		logger.Println(file.lines[row])

		// print an arrow with a tail spanning the width of the mistake
		for column > 0 {
			logger.Print(" ")
			column --
		}
		for width > 1 {
			logger.Print("-")
		}
		logger.Println("^")
	}
	logger.Println(message)
}

// Error prints an error at row and column spanning width amount of runes, and
// a message describing the error. If width is zero, neither the line nor the
// location information is printed.
func (file *File) Error (column, row, width int, message string) {
	file.mistake(column, row, width, message, "ERR")
}

// Warn prints a warning at row and column spanning width amount of runes, and
// a message describing the error. If width is zero, neither the line nor the
// location information is printed.
func (file *File) Warn (column, row, width int, message string) {
	file.mistake(column, row, width, message, "!!!")
}

// Location returns a location struct describing the current position inside of
// the file. This can be stored and used to print errors.
func (file *File) Location () (location Location) {
	return Location {
		file:   file,
		row:    file.currentLine,
		column: file.currentColumn,
	}
}
