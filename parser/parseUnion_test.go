package parser

import (
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestUnion1(t *testing.T) {

	input := `
union SearchResult =| Photo | Person

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
	expDoc := `type Person {name:String age:Int} type Photo { height:Intwidth:Int} type SearchQuery{firstSearchResult:SearchResult} union SearchResult=|Photo|Person`

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
	if compare(d.String(), expDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestUnionDupMember(t *testing.T) {

	input := `
union SearchResult =| Photo | Person | Photo
`

	var expectedErr [1]string
	expectedErr[0] = `Duplicate member name at line: 2 column: 40`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

}

func TestUnionInvalidMember(t *testing.T) {

	input := `
input OddOne {
	x: Int
	y: Float
}
union SearchResult =| Photo | Person | OddOne
`

	var expectedErr [1]string
	expectedErr[0] = `Union member "OddOne" must be an object based type at line: 6 column: 40`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

}
