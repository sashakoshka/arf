package analyzer

import "os"
import "fmt"
import "path/filepath"
// import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// AnalysisOperation holds information about an ongoing analysis operation.
type AnalysisOperation struct {
	sectionTable SectionTable
	modulePath   string

	currentPosition locator
	currentSection  parser.Section
	currentTree     parser.SyntaxTree
}

// Analyze performs a semantic analyisys on the module specified by path, and
// returns a SectionTable that can be translated into C.
func Analyze (modulePath string, skim bool) (table SectionTable, err error) {
	if modulePath[0] != '/' {
		cwd, _ := os.Getwd()
		modulePath = filepath.Join(cwd, modulePath)
	}

	analyzer := AnalysisOperation {
		sectionTable: make(SectionTable),
		modulePath:   modulePath,
	}

	err = analyzer.analyze()
	table = analyzer.sectionTable
	return
}

// analyze performs an analysis operation given the state of the operation
// struct.
func (analyzer *AnalysisOperation) analyze () (err error) {
	tree, err := parser.Fetch(analyzer.modulePath, false)
	sections := tree.Sections()

	for !sections.End() {
		_, err = analyzer.fetchSection(locator {
			modulePath: analyzer.modulePath,
			name: sections.Value().Name(),
		})
		sections.Next()
	}
	
	return
}

// fetchSection returns a section from the section table. If it has not already
// been analyzed, it analyzes it first. If the section does not actually exist,
// a nil section is returned. When this happens, an error should be created on
// whatever syntax tree node "requested" the section be analyzed.
func (analyzer *AnalysisOperation) fetchSection (
	where locator,
) (
	section Section,
	err error,
) {
	var exists bool
	section, exists = analyzer.sectionTable[where]
	if exists { return }

	// fetch the module. since we already have our main module parsed fully
	// and not skimmed, we can just say "yeah lets skim stuff here".
	var tree parser.SyntaxTree
	tree, err = parser.Fetch(where.modulePath, true)
	if err != nil {
		section = nil
		return
	}

	var parsedSection = tree.LookupSection(where.name)
	if parsedSection == nil {
		section = nil
		return
	}

	previousPosition := analyzer.currentPosition
	previousSection  := analyzer.currentSection
	previousTree     := analyzer.currentTree
	analyzer.currentPosition = where
	analyzer.currentSection  = parsedSection
	analyzer.currentTree     = tree

	defer func () {
		analyzer.currentPosition = previousPosition
		analyzer.currentSection  = previousSection
		analyzer.currentTree     = previousTree
	} ()

	// TODO: analyze section. have analysis methods work on currentPosition
	// and currentSection.
	// 
	// while building an analyzed section, add it to the section
	// table as soon as the vital details are acquired, and mark it as
	// incomplete. that way, it can still be referenced by itself in certain
	// scenarios.
	switch parsedSection.Kind() {
	case parser.SectionKindType:
		section, err = analyzer.analyzeTypeSection()
		if err != nil { return}
	case parser.SectionKindEnum:
	case parser.SectionKindFace:
	case parser.SectionKindData:
	case parser.SectionKindFunc:
	}
	
	return
}

// fetchSectionFromIdentifier is like fetchSection, but takes in an identifier
// referring to a section and returns the section. This works within the context
// of whatever module is currently being analyzed. The identifier in question
// may have more items than 1 or 2, but those will be ignored. This method
// "consumes" items from the identifier, it will return an identifier without
// those items.
func (analyzer *AnalysisOperation) fetchSectionFromIdentifier (
	which parser.Identifier,
) (
	section Section,
	bitten  parser.Identifier,
	err     error,
) {
	bitten = which
	item := bitten.Bite()
	
	path, exists := analyzer.currentTree.ResolveRequire(item)
	if exists {
		// we have our module path, so get the section name
		item = bitten.Bite()
	} else {
		// that wasn't a module name, so the module path must be the our
		// current one
		path = analyzer.currentPosition.modulePath
	}

	section, err = analyzer.fetchSection (locator {
		name:       item,
		modulePath: path,
	})
	if err != nil { return }

	if section == nil {
		err = which.NewError (
			"section \"" + item + "\" does not exist",
			infoerr.ErrorKindError,
		)
		return
	}
	
	return
}

func doIndent (indent int, input ...any) (output string) {
	for index := 0; index < indent; index ++ {
		output += "\t"
	}

	output += fmt.Sprint(input...)
	return
}
