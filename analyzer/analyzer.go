package analyzer

import "os"
import "path/filepath"
import "git.tebibyte.media/arf/arf/parser"

// AnalysisOperation holds information about an ongoing analysis operation.
type AnalysisOperation struct {
	sectionTable SectionTable
	modulePath   string
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
		switch sections.Value().Kind() {
			
		}
		
		sections.Next()
	}
	
	return
}
