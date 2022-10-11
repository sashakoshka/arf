package analyzer

import "os"
import "sort"
import "path/filepath"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

// locator uniquely identifies a section in the section table.
type locator struct {
	modulePath string
	name string
}

func (where locator) ToString () (output string) {
	cwd, _ := os.Getwd()
	modulePath, err := filepath.Rel(cwd, where.modulePath)
	if err != nil {
		panic("cant get relative path: " + err.Error())
	}

	output += modulePath + "." + where.name
	return
}

// SectionTable stores a list of semantically analized sections from one module,
// and all sections that it requires from other modules.
type SectionTable map[locator] Section

// ToString returns the data stored in the table as a string.
func (table SectionTable) ToString (indent int) (output string) {
	sortedKeys := make(locatorArray, len(table))
	index := 0
	for key, _ := range table {
		sortedKeys[index] = key
		index ++
	}
	sort.Sort(sortedKeys)
	
	for _, name := range sortedKeys {
		section := table[name]
		output += section.ToString(indent)
	}

	return
}

// locatorArray holds a sortable array of locators
type locatorArray []locator

// Len returns the length of the locator array
func (array locatorArray) Len () (length int) {
	length = len(array)
	return
}

// Less returns whether item at index left is less than item at index right.
func (array locatorArray) Less (left, right int) (less bool) {
	leftLocator  := array[left]
	rightLocator := array[right]
	
	less =
		leftLocator.modulePath  + leftLocator.name <
		rightLocator.modulePath + rightLocator.name
	return
}

// Swap swaps the elments at indices left and right.
func (array locatorArray) Swap (left, right int) {
	temp := array[left]
	array[left]  = array[right]
	array[right] = temp
}

// Section is a semantically analyzed section.
type Section interface {
	// Provided by sectionBase
	Name       () (name string)
	Complete   () (complete bool)
	ModulePath () (path string)
	ModuleName () (path string)
	Permission () (permission types.Permission)
	Location   () (location file.Location)
	NewError   (message string, kind infoerr.ErrorKind) (err error)
	locator    () (where locator)

	// Must be implemented by each individual section
	ToString (indent int) (output string)
}
