/*
Package analyzer implements a semantic analyzer for the ARF language. In it,
there is a function called Analyze which takes in a module path and returns a
table of sections representative of that module. The result of this operation is
not cached.

The section table returned by Analyze can be expected to be both syntactically
correct and semantically sound.

This package automatically invokes the parser and lexer packages.
*/
package analyzer

import "os"
import "fmt"
import "path/filepath"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/parser"
import "git.tebibyte.media/arf/arf/infoerr"

// analysisOperation holds information about an ongoing analysis operation.
type analysisOperation struct {
	sectionTable SectionTable
	modulePath   string

	currentPosition locator
	currentSection  parser.Section
	currentTree     parser.SyntaxTree
}

// Analyze performs a semantic analysis on the module specified by path, and
// returns a SectionTable that can be translated into C. The result of this is
// not cached.
func Analyze (modulePath string, skim bool) (table SectionTable, err error) {
	if modulePath[0] != '/' {
		cwd, _ := os.Getwd()
		modulePath = filepath.Join(cwd, modulePath)
	}

	analyzer := analysisOperation {
		sectionTable: make(SectionTable),
		modulePath:   modulePath,
	}

	err = analyzer.analyze()
	table = analyzer.sectionTable
	return
}

// analyze performs an analysis operation given the state of the operation
// struct.
func (analyzer *analysisOperation) analyze () (err error) {
	var tree parser.SyntaxTree
	tree, err = parser.Fetch(analyzer.modulePath, false)
	if err != nil { return }
	sections := tree.Sections()

	for !sections.End() {
		_, err = analyzer.fetchSection(locator {
			modulePath: analyzer.modulePath,
			name: sections.Value().Name(),
		})
		if err != nil { return err }
		sections.Next()
	}
	
	return
}

// fetchSection returns a section from the section table. If it has not already
// been analyzed, it analyzes it first. If the section does not actually exist,
// a nil section is returned. When this happens, an error should be created on
// whatever syntax tree node "requested" the section be analyzed.
func (analyzer *analysisOperation) fetchSection (
	where locator,
) (
	section Section,
	err error,
) {
	var exists bool
	section, exists = analyzer.resolvePrimitive(where)
	if exists { return }

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

	// analyze section. have analysis methods work on currentPosition
	// and currentSection.
	// 
	// while building an analyzed section, add it to the section
	// table as soon as the vital details are acquired, and mark it as
	// incomplete. that way, it can still be referenced by itself in certain
	// scenarios.
	switch parsedSection.(type) {
	case parser.TypeSection:
		section, err = analyzer.analyzeTypeSection()
		if err != nil { return}
	case parser.EnumSection:
		section, err = analyzer.analyzeEnumSection()
		if err != nil { return}
	case parser.FaceSection:
	case parser.DataSection:
		section, err = analyzer.analyzeDataSection()
		if err != nil { return}
	case parser.FuncSection:
	}
	
	return
}

// TODO: make this method a generalized "get this from an identifier in context
// of the current scope" method. have it return various things like sections,
// variables, functions, members, methods, etc. have it return an any.
// if no node could be found, return an error saying entity not found. if a node
// is found, but other identifier items describe members of that node that do
// not exist, return an error saying nonexistent member. if the node is private,
// also return an error.
//
// when new things are defined, they should not be allowed to shadow anything
// else in above scopes. nevertheless, the method should search in this order:
// 
// 1. search scopes starting with closest -> farthest
// 2. if first part of identifier is a require, get section from other module
// 3. search for section in current module
//
// look into making a unified structure for data sections and variables, and
// having data section variables be part of a "root" scope at the base of every
// module.

// fetchNodeFromIdentifier is like fetchSection, but takes in an identifier
// referring to any node accessible within the current scope and returns it.
// This method works within the current scope and current module. This method
// consumes the entire identifier, and will produce an error if there are
// identifier items left unconsumed.
func (analyzer *analysisOperation) fetchNodeFromIdentifier (
	which parser.Identifier,
) (
	node any,
	err  error,
) {
	item, bitten := which.Bite()

	// TODO: search scopes for variables

	// the identifier must be referring to a section
	var external bool
	path, exists := analyzer.currentTree.ResolveRequire(item)
	if exists {
		// we have our module path, so get the section name
		item, bitten = bitten.Bite()
		external = true
	} else {
		// that wasn't a module name, so the module path must be our
		// current one
		path = analyzer.currentPosition.modulePath
	}

	// attempt to get section
	var section Section
	section, err = analyzer.fetchSection (locator {
		name:       item,
		modulePath: path,
	})
	node = section
	if err != nil { return }

	// return error if nothing mentioned in the identifier is accessible
	if node == nil {
		err = which.NewError (
			"can't find anything called \"" + item + "\" within " +
			"current scope",
			infoerr.ErrorKindError,
		)
		return
	}

	// return error if the section is private
	if external && section.Permission() == types.PermissionPrivate {
		err = which.NewError(
			"this section is private, and cannot be used " +
			"outside of its module",
			infoerr.ErrorKindError)
		return
	}

	// see if we can do member selection on the section
	// TODO: at this point, we are gonna return an argument.
	if bitten.Length > 0 {
		switch 
	}
	
	return
}

// resolvePrimitive checks to see if the locator is in the current module, and
// refers to a primitive. If it does, it returns a pointer to that primitive
// and true for exists. If it doesn't, it returns nil and false.
func (analyzer *analysisOperation) resolvePrimitive (
	where locator,
) (
	section Section,
	exists bool,
) {
	// primitives are scoped as if they are contained within the current
	// module, so if the location refers to something outside of the current
	// module, it is definetly not referring to a primitive.
	if where.modulePath != analyzer.currentPosition.modulePath {
		return
	}

	exists = true
	switch where.name {
	case "Int":    section = &PrimitiveInt
	case "UInt":   section = &PrimitiveUInt
	case "I8":     section = &PrimitiveI8
	case "I16":    section = &PrimitiveI16
	case "I32":    section = &PrimitiveI32
	case "I64":    section = &PrimitiveI64
	case "U8":     section = &PrimitiveU8
	case "U16":    section = &PrimitiveU16
	case "U32":    section = &PrimitiveU32
	case "U64":    section = &PrimitiveU64
	case "Obj":    section = &PrimitiveObj
	// case "Face":   section = &PrimitiveFace
	// case "Func":   section = &PrimitiveFunc
	case "String": section = &BuiltInString
	default:
		exists = false
	}

	return
}

// addSection adds a section to the analyzer's section table. If a section with
// that name already exists, it panics because the parser should not have given
// that to us.
func (analyzer *analysisOperation) addSection (section Section) {
	_, exists := analyzer.sectionTable[section.locator()]
	if exists {
		panic (
			"invalid state: duplicate section " +
			section.locator().ToString())
	}
	analyzer.sectionTable[section.locator()] = section
	return
}

// typeCheck checks to see if source can fit as an argument into a slot of type
// destination. If it can, it retunrs nil. If it can't, it returns an error
// explaining why.
func (analyzer *analysisOperation) typeCheck (
	source Argument,
	destination Type,
) (
	err error,
) {
	if !source.canBePassedAs(destination) {
		err = source.NewError (
			typeMismatchErrorMessage (
				source.What(),
				destination),
			infoerr.ErrorKindError)
	}

	return
}

// inCurrentModule returns whether or not the specified section resides within
// the current module.
func (analyzer *analysisOperation) inCurrentModule (
	section Section,
) (
	inCurrent bool,
){
	inCurrent =
		section.locator().modulePath ==
		analyzer.currentPosition.modulePath
	return
}

// TODO: make a method of analyzer that, given a name, searches through all
// accessible scopes and returns the thing the name references. when analyzing
// a function, the analyzer should remember a trail of scopes.

// doIndent perfroms a fmt.Sprint operation on input, indenting the string. This
// does not add a trailing newline.
func doIndent (indent int, input ...any) (output string) {
	for index := 0; index < indent; index ++ {
		output += "\t"
	}

	output += fmt.Sprint(input...)
	return
}
