package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/db"
	"github.com/graph-sdl/lexer"
)

func TestSetup4Fragments(t *testing.T) {

	input := `
enum LengthUnit {
    METER
    CENTERMETER
    MILLIMETER
    KILOMETER
}

enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
  DRTYPE
}

type Starship {
  id: ID!
  name: String!
  length(unit: LengthUnit = METER): Float
}

interface Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode!]!
}

type Human implements Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode!]!
  starships: [Starship]
  totalCredits: Int
}

type Droid implements Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode!]!
  primaryFunction: String
}

type Query {
  hero(episode: Episode): [Character]
  droid(id: ID!): Droid
}


`

	expectedDoc := `interface Character {id:ID! name:String! friends:[Character] appearsIn:[Episode!]!}
	
	 type Droid implements,Character{id:ID! name:String! friends:[Character] appearsIn:[Episode!]! primaryFunction:String}
	 
	 enum Episode{NEWHOPE EMPIRE JEDI DRTYPE}
	 
	type Human implements Character{id:ID! name:String! friends:[Character] appearsIn:[Episode!]! starships:[Starship] totalCredits:Int}

	enum LengthUnit{METER CENTERMETER MILLIMETER KILOMETER}
 	type Query {
  hero(episode: Episode): [Character]
  droid(id: ID!): Droid
}   
    type Starship {id:ID! name:String! length(unit:LengthUnit=METER):Float}
   `
	var expectedErr [3]string
	expectedErr[0] = ``
	// 	expectedErr[1] = `Type "Business" does not implement interface "NamedEntity", missing  "XXX"`
	// 	expectedErr[2] = `Type "Business" does not implement interface "ValuedEntity", missing  "size" "length"`

	err := db.DeleteType("Starship")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Character")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Human")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Droid")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println("*** Error: ", v)
	}
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestSetup4Union(t *testing.T) {

	input := `

union SearchResult = Photo | Person

type Person {
  name: String
  age: Int
}

type Photo {
  height: Int
  width: Int
}

type SearchQuery {
  firstSearchResult: SearchResult
}

type Query {
  hero(episode: Episode = JEDI ): SearchQuery
}

enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
  DRTYPE
}

`

	expectedDoc := `
enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
  DRTYPE
}
type Person {
  name: String
  age: Int
}

type Photo {
  height: Int
  width: Int
}
type Query {
  hero(episode: Episode = JEDI): SearchQuery
}

type SearchQuery {
  firstSearchResult: SearchResult
}
union SearchResult = |Photo | Person
`

	// 	var expectedErr [3]string

	err := db.DeleteType("SearchQuery")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Photo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SearchResult")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestSetup4UnionWithDefaultError(t *testing.T) {

	input := `

union SearchResult = Photo | Person

type Person {
  name: String
  age: Int
}

type Photo {
  height: Int
  width: Int
}

type SearchQuery {
  firstSearchResult: SearchResult
}

type Query {
  hero(episode: Episode = JEDII ): SearchQuery
}

enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
  DRTYPE
}

`

	expectedDoc := `
enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
  DRTYPE
}
type Person {
  name: String
  age: Int
}

type Photo {
  height: Int
  width: Int
}
type Query {
  hero(episode: Episode = JEDII): SearchQuery
}

type SearchQuery {
  firstSearchResult: SearchResult
}
union SearchResult = |Photo | Person
`

	var expectedErr [1]string
	expectedErr[0] = `"JEDII" is not a member of Enum type Episode at line: 20 column: 27`

	err := db.DeleteType("SearchQuery")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Photo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SearchResult")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}
