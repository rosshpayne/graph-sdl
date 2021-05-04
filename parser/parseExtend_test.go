package parser

import (
	"testing"

	"github.com/graphql/internal/graph-sdl/lexer"
)

func TestExtend0(t *testing.T) {

	input := `
input ExampleInputObject {
  a: String
  b: Int!
}


`

	//expectedDoc := `input ExampleInputObject {a : String b : Int!  }`

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

func TestExtend1(t *testing.T) {

	input := `


extend input ExampleInputObject  {
  name: String
  age: [Int!]
}

`

	expectedDoc := `input ExampleInputObject {a : String b : Int! name : String age : [Int!] }`

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

// func TestExtend1a(t *testing.T) {

// 	input := `

// extend input Address  {
//   name: String
//   age: [Int!]
// }

// `

// 	expectedDoc := `input ExampleInputObject {a : String b : Int! name : String age : [Int!] }`

// 	l := lexer.New(input)
// 	p := New(l)
// 	d, errs := p.ParseDocument()
// 	fmt.Println(d.String())
// 	for i, v := range errs {
// 		fmt.Println(i, v.Error())
// 	}
// 	if compare(d.String(), expectedDoc) {
// 		fmt.Println(trimWS(d.String()))
// 		fmt.Println(trimWS(expectedDoc))
// 		t.Errorf(`*************  program.String() wrong.`)
// 	}
// }

func TestExtendDupField(t *testing.T) {

	input := `

extend input ExampleInputObject {
	age: [Int!]
}

`

	var expectedErr [1]string
	expectedErr[0] = `Duplicate input value name "age" at line: 4, column: 2`

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
