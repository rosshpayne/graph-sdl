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

	expectedErr := `Illegal IDENT token, [+i] at line: 3, column: 12`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) != 1 {
		t.Errorf(`Not expected - should be 1 error got %d`, len(errs))
	}
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s] `, v.Error(), expectedErr)
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
		t.Errorf(`*************  program.String() wrong.`)
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

	expectedErr := `Expected name identifer got NULL of "null" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
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

	expectedErr := `Expected a ( or } or { instead got [ at line: 5, column: 26`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
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

	expectedErr := `Expected name identifer got TRUE of "true" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}

func TestMutation1(t *testing.T) {

	input := `type Mutation {
  createPerson(name: String!, age: Int!): [[PersonX!]]
}
`

	expectedErr := `Type "PersonX", not defined at line: 2 column: 45`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	//fmt.Println(d.String())
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
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

	var expectedErr [4]string
	expectedErr[0] = `Expected name identifer got TRUE of "true" at line: 4, column: 3`
	expectedErr[1] = `identifer [__dep] cannot start with two underscores at line: 6, column: 22`
	expectedErr[2] = `Type "Place", not defined at line: 9 column: 13`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for i, v := range errs {
		fmt.Println(i, v.Error())
	}
	if len(errs) != 3 {
		t.Errorf(`Expect 4 errors got %d`, len(errs))
	} else {
		for i, v := range errs {
			if i < 3 && v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		}
	}
}
