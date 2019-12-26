package parser

import (
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
	 
	 enum Episode{NEWHOPE EMPIRE JEDI}
	 
	type Human implements Character{id:ID! name:String! friends:[Character] appearsIn:[Episode!]! starships:[Starship] totalCredits:Int}

	enum LengthUnit{METER CENTERMETER MILLIMETER KILOMETER}
 	type Query {
  hero(episode: Episode): [Character]
  droid(id: ID!): Droid
}   
    type Starship {id:ID! name:String! length(unit:LengthUnit=METER):Float}
   `
	// 	var expectedErr [3]string
	// 	expectedErr[0] = `Type "Person" does not implement interface "NamedEntity", missing  "XXX"`
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

`

	expectedDoc := `

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
union SearchResult = |Photo | Person

   `
	// 	var expectedErr [3]string
	// 	expectedErr[0] = `Type "Person" does not implement interface "NamedEntity", missing  "XXX"`
	// 	expectedErr[1] = `Type "Business" does not implement interface "NamedEntity", missing  "XXX"`
	// 	expectedErr[2] = `Type "Business" does not implement interface "ValuedEntity", missing  "size" "length"`

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
