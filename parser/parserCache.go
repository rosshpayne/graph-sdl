package parser

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/graph-sdl/ast"
	"github.com/graph-sdl/db"
	"github.com/graph-sdl/lexer"
)

// for GL types only
type entry struct {
	ready chan struct{}       // a channel for each entry - to synchronise access when the data is being sourced
	data  ast.GQLTypeProvider // this represents the AST data to be cached. Its value is populated after the entry is saved in the cache.
}

type Cache_ struct {
	sync.Mutex // Mutex protects whole cache. Channels protect individual cache entries.
	Cache      map[string]*entry
	logr       *log.Logger
}

func (tc *Cache_) SetLogger(logr *log.Logger) {
	tc.logr = logr
}

// NewCache allocates a structure to hold the cached data with access methods.
func NewCache() *Cache_ {
	typeNotExists = make(map[string]bool)
	return &Cache_{Cache: make(map[string]*entry)}
}

// AddEntry is not concurrency safe. Used in non-cconcurrent situations.
func (t *Cache_) AddEntry(name ast.NameValue_, data ast.GQLTypeProvider) { //ast.NameValue_, data GQLTypeProvider) {
	e := &entry{data: data, ready: make(chan struct{})}
	close(e.ready)
	fmt.Println("Added to cache ", name.String())
	t.Cache[name.String()] = e
}

var (
	typeNotExists map[string]bool

	// errors
	ErrNotCached error = errors.New("does not exist")
	ErrnotScalar error = errors.New("scalars are not permitted for FetchAST")
	ErrnoName    error = errors.New("No input name supplied to FetchAST")
)

// FetchAST is a concurrency safe access method to the cache. If entry not found in the cache searches dynamodb table for the type SDL statement.
//  Parses statement to create the AST and populates the cache and returns the AST.
func (t *Cache_) FetchAST(name ast.NameValue_) (ast.GQLTypeProvider, error) {

	fmt.Println("**** FetchAST ", name.String())
	name_ := name.String()
	//
	// do not handle scalars or nul name
	switch name_ {
	case "String", "Int", "Float", "Boolean", "ID", "null":
		return nil, ErrnotScalar
	}
	if len(name) == 0 {
		return nil, ErrnoName
	}
	// check if name has been registered as non-existent from previous query
	if typeNotExists[name_] {
		fmt.Printf("DBFetch of [%s] does not exist", name)
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
			if errors.Is(err, db.SystemErr) {
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
				t.logr.Printf("Found in DB: %q\n", typeSDL)
				// generate AST for the resolved type
				t.logr.Print(" in parseCache about to generate AST.")
				l := lexer.New(typeSDL)
				p2 := New(l)
				e.data = p2.ParseStatement() // source of stmt is db so its been verified, simply resolve types it refs
				p2.ResolveAllTypes(e.data, t)
				// close the channel to allow unhindered access to this entry
				t.logr.Print(" found in db. closed channel...")
				close(e.ready)
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
	fmt.Println("**** FetchAST returned with data ", e.data.TypeName())
	return e.data, nil

}

func (t *Cache_) CacheClear() {
	fmt.Println("******************************************")
	fmt.Println("************ CLEAR CACHE *****************")
	fmt.Println("******************************************")
	t.Cache = make(map[string]*entry)
}
