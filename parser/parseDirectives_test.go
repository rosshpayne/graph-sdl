package parser

import (
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestMultiDirective1(t *testing.T) {

	input := `
input ExampleInputObjectDirective @ june (asdf:234) @ june2 (aesdf:234) @ june3 (as2df:"abc") {
  a: String = "AbcDef" @ ref (if:123) @ jack (sd: "abc") @ june (asdf:234) @ ju (asdf:234) @ judkne (asdf:234) @ junse (asdf:234) @ junqe (asdf:234) 
  b: Int!@june(asdf:234) @ ju (asdf:234)
}
`

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
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}

}

func TestInputDoesnotExist(t *testing.T) {

	input := `
extend input ExampleInputXYZ @ june (asdf:234) 
`
	var expectedErr [1]string
	expectedErr[0] = `Type "ExampleInputXYZ" does not exist at line: 2 column: 14`

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

func TestExtendInpDirDuplicate(t *testing.T) {

	input := `
extend input ExampleInputObjectDirective @ june (asdf:234) 
`
	var expectedErr [2]string
	expectedErr[0] = `Duplicate Directive name "june" at line: 2, column: 44`
	expectedErr[1] = `extend for type "ExampleInputObjectDirective" contains no changes at line: 0, column: 0`

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
