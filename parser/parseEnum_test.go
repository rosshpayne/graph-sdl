package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestBadTypeName(t *testing.T) {

	input := `
type Person {
  name: Str+ing!
  address: String
  Altitude: Float
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got ILLEGAL of "+" at line: 3, column: 12`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be 2 errors got %d`, len(errs)))
	}
	//fmt.Println(d.String())

	for _, got := range errs {
		var found bool
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Not expected Error =[%q]`, got.Error())
		}
	}
}

func TestEnum1X(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
}
`

	l := lexer.New(input)
	p := New(l)
	d, _ := p.ParseDocument()
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		if len(trimWS(input)) != len(trimWS(d.String())) {
			t.Errorf(`*************  program.String() wrong.`)
		}
	}

}

func TestBadEnumValue1(t *testing.T) {

	input := `
enum Direction {
  NORTH
  null
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got null of "null" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestBadBracket(t *testing.T) {

	input := `
	enum Direction {
  NORTH
  SOUTH
  WEST @deprecated @ dep [if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected a ( or } or { instead got [ at line: 5, column: 26`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestBadEnumValue2(t *testing.T) {

	input := `
	enum Direction {
  NORTH
  true
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got TRUE of "true" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestMutation1(t *testing.T) {

	input := `type Mutation {
  createPerson(name: String!, age: Int!): [[PersonX]]
}
`
	var expectedErr [1]string
	expectedErr[0] = `Type "PersonX" does not exist at line: 2 column: 45`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestMultiError1(t *testing.T) {

	input := `
enum Direction {
  NORTH
  true
  SOUTH
  WEST @deprecated @ __dep (if: 99.34)
}
type Person {
		address: [Place]
		name: [String!]!
	}
`

	var expectedErr [3]string
	expectedErr[0] = `Expected name identifer got TRUE of "true" at line: 4, column: 3`
	expectedErr[1] = `identifer [__dep] cannot start with two underscores at line: 6, column: 22`
	expectedErr[2] = `Type "Place" does not exist at line: 9 column: 13`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}
