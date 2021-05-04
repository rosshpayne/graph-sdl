package parser

import (
	"errors"
	"log"
	"sync"

	"github.com/graphql/internal/graph-sdl/ast"
	"github.com/graphql/internal/graph-sdl/db"
	"github.com/graphql/internal/graph-sdl/lexer"
)

// for GL types only
type entry struct {
	ready chan struct{}       // a channel for each entry - to synchronise access when the data is being sourced
	data  ast.GQLTypeProvider // this represents the AST data to be cached. Its value is populated after the entry is saved in the cache.
}

type Cache_ struct {
	sync.Mutex                   // Mutex protects whole cache. Channels protect individual cache entries.
	Cache      map[string]*entry // cache holds any AST accessed by its name such as types, statements in a doc.
	logr       *log.Logger
}

// instance of a cache. This is shared amoungst all parser and query executers.
var cache *Cache_

func (tc *Cache_) SetLogger(logr *log.Logger) {
	if tc.logr == nil {
		tc.logr = logr
	}
}

// init creates two caches, the not-exists cache which contain all types that do not exist in the current document or in the document being parsed.
// The other cache is for all types that are being created from the parsed document and those that exist in the database. It is populated as required.
// The cache exists at the package level, so is available to each parser. The alternate design is to not use init and create the caches in NewCache() below.
func init() {
	typeNotExists = make(map[string]bool)
	cache = &Cache_{Cache: make(map[string]*entry)}
}

// NewCache allocates a structure to hold the cached data with access methods.
func NewCache() *Cache_ {
	return cache
	//	typeNotExists = make(map[string]bool)
	//	return &Cache_{Cache: make(map[string]*entry)} // note: this design has each parser/executer assigned its own cache. No concurrency issues but requires more memory
	//  and one parser/executor doesn't benefit from the work of others. Also more db IO.

}

// AddEntry to cache is  concurrency safe.
// TODO: check typeNotExists cache is handled safely. Concurrency was designed around Cache not typeNotExists cache.
// NOTE: AddEntry Does NOT save to database. The type is saved to the db only if it has zero errors as part of the defer in parser.go
func (t *Cache_) addEntry(name ast.NameValue_, data ast.GQLTypeProvider) { //ast.NameValue_, data GQLTypeProvider) {
	e := &entry{data: data, ready: make(chan struct{})}
	close(e.ready)
	t.Lock()
	// delete from notExists cache - if present
	delete(typeNotExists, name.String())
	// add to type cache
	t.Cache[name.String()] = e
	t.Unlock()
	t.logr.Println("addEntry:  Added to Type cache ", name)
}

var (
	typeNotExists map[string]bool

	// errors
	ErrNotCached error = errors.New("does not exist")
	ErrnotScalar error = errors.New("scalars are not permitted for FetchAST")
	ErrnoName    error = errors.New("No input name supplied to FetchAST")
)

// FetchAST is a concurrency safe access method to the cache. Used when resolving nested abstract types for the type being created.
// When all validation checks are satisfieid the type in question is added to the cache.
// If entry not found in the cache searches dynamodb table for the type SDL statement.
func (t *Cache_) FetchAST(name ast.NameValue_) (ast.GQLTypeProvider, error) {

	name_ := name.String()
	//
	t.logr.Println("FetchAST: ", name)
	// do not handle scalars or nul name
	switch name_ {
	case "String", "Int", "Float", "Boolean", "ID", "null":
		t.logr.Println("FetchAST: scalar value abort...")
		return nil, ErrnotScalar
	}
	if len(name) == 0 {
		return nil, ErrnoName
	}
	// check if name has been registered as non-existent from previous query
	if typeNotExists[name_] {
		t.logr.Printf("FetchAST: DBFetch of [%s] does not exist\n", name)
		return nil, ErrNotCached
	}
	t.Lock()
	e := t.Cache[name_] // e will be nil only when name_ is not in the cache. Nil has no other meaning.

	if e == nil {

		e = &entry{ready: make(chan struct{})}
		// save pointer entry struct to cache now. The AST struct field will be assigned to struct soon. Channel comms will comunicate when AST is populated
		t.Cache[name_] = e
		t.Unlock()
		// cache populated with bare minimum of data.  Release the lock and source remaining data to be cached while the channel synchronises access to the current entry.
		// access db for definition of type (string value)
		if typeSDL, err := db.DBFetch(name_); err != nil {
			switch {
			case errors.Is(err, db.SystemErr), errors.Is(err, db.MarshalingErr), errors.Is(err, db.UnmarshalingErr):
				t.logr.Fatal(err)
			}
			typeNotExists[name_] = true
			delete(t.Cache, name_)
			close(e.ready)
			if errors.Is(err, db.NoItemFoundErr) {
				t.logr.Print(err)
			}
			return nil, err
		} else {
			if len(typeSDL) == 0 { // no type found in DB
				// mark type as being nonexistent
				t.logr.Print("Type not found ")
				typeNotExists[name_] = true
				delete(t.Cache, name_)
				close(e.ready)
				return nil, err
			} else {
				t.logr.Printf("FetchAST:  returned from DB: %q\n", typeSDL)
				// generate AST for the resolved type
				l := lexer.New(typeSDL)
				p2 := New(l)
				//
				// Generate AST for name of stmt or a GQL type and save to cache
				// Important: source of stmt is db so its been verified, simply resolve types it refs
				ast_ := p2.ParseStatement()
				e.data = ast_
				close(e.ready)
				//
				// resolve dependent types
				//
				p2.resolveDependents(ast_, t)
			}
		}
	} else {
		t.Unlock()
		<-e.ready // AST is now populated in cache for this named type
	}
	if e.data == nil {
		// concurrency issue  (when currency applies) - two queries on same object within short time interval - before typeNotExists is updated.
		return nil, ErrNotCached
	}
	return e.data, nil

}

func (t *Cache_) CacheClear() {
	if t == nil {
		return
	}
	t.Lock()
	t.Cache = make(map[string]*entry)
	typeNotExists = make(map[string]bool)
	t.Unlock()
}
