package parser

// cacheItem stores an item of the parser cache.
type cacheItem struct {
	tree    SyntaxTree
	skimmed bool
}

// cache stores all modules that have been parsed so far. They are indexed with
// their full path on the filesystem, starting with '/'.
var cache = make(map[string] cacheItem)
