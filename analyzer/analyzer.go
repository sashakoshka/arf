package analyzer

import "os"
import "fmt"
import "path/filepath"
// import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
// import "git.tebibyte.media/arf/arf/infoerr"

// AnalysisOperation holds information about an ongoing analysis operation.
type AnalysisOperation struct {
	sectionTable SectionTable
	modulePath   string

	currentPosition locator
	currentSection  parser.Section
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
	analyzer.currentPosition = where
	analyzer.currentSection  = parsedSection

	// TODO: analyze section. have analysis methods work on currentPosition
	// and currentSection.
	// 
	// while building an analyzed section, add it to the section
	// table as soon as the vital details are acquired, and mark it as
	// incomplete. that way, it can still be referenced by itself in certain
	// scenarios.
	switch parsedSection.Kind() {
	case parser.SectionKindType:
	case parser.SectionKindObjt:
	case parser.SectionKindEnum:
	case parser.SectionKindFace:
	case parser.SectionKindData:
	case parser.SectionKindFunc:
	}
	
	analyzer.currentPosition = previousPosition
	analyzer.currentSection  = previousSection
	return
}

func doIndent (indent int, input ...any) (output string) {
	for index := 0; index < indent; index ++ {
		output += "\t"
	}

	output += fmt.Sprint(input...)
	return
}
