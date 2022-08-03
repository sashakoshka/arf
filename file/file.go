package file

import "os"
import "path/filepath"

type File struct {
	path        string
	moduleName  string
	file        *os.File
	currentLine int
	lines       [][]byte
}

// Open opens the file specified by path and returns a new File struct. The
// module name is automatically inferred through hte file's parent directory.
func Open (path string) (file *File, err error) {
	file = &File {
		path:       path,
		moduleName: filepath.Dir(path),
		lines:      [][]byte { []byte { } },
	}

	file.file, err = os.OpenFile(path, os.O_RDONLY, 0660)
	if err != nil { return }

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
	amountRead, err = file.file.Read(bytes)

	for _, char := range bytes {
		if char == '\n' {
			file.lines = append(file.lines, []byte { })
			file.currentLine ++
		} else {
			file.lines[file.currentLine] = append (
				file.lines[file.currentLine],
				bytes...)
		}
	}
	
	return
}

// Close closes the file. After the file is closed, data that has been read will
// still be retained, and errors can be reported.
func (file *File) Close () {
	file.file.Close()
}
